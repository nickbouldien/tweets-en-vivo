package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

const (
	ApiToken = "API_TOKEN"
	baseURL = "https://api.twitter.com/2"
	rulesURL = baseURL + "/tweets/search/stream/rules"
	streamURL = baseURL + "/tweets/search/stream"
)

var apiToken string

type DeleteIDs []string

type DeleteRules struct {
	Delete map[string]DeleteIDs `json:"delete"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(`Error loading the .env file. make sure it exists and 
				has all of the required environment variables`)
	}

	command := flag.String("command", "stream", "the command you want to perform")
	file := flag.String("file", "rules.json", "the rules file you want to use (in the rules/ dir)")
	flag.Parse()

	fmt.Println("file:", *file)
	fmt.Println("command:", *command)

	//if flag.NArg() < 1 {
	//	fmt.Printf("usage:\n\t%s \"tweets-en-vivo\"\n", os.Args[0])
	//	os.Exit(1)
	//}

	apiToken = os.Getenv(ApiToken)

	if apiToken == "" {
		log.Fatal(`make sure that you have filled in the required
				twitter api credentials in the .env file`)
	}

	switch *command {
	case "add":
		fmt.Println("add")
		// import the rules json file
		jsonFile, err := os.Open(path.Join("rules/", *file))
		if err != nil {
			log.Fatal("could not open the json file", err)
		}
		defer jsonFile.Close()

		byteValue, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			log.Fatal("could not read the json", err)
		}

		// add the rules
		body, err := AddRules(byteValue, false)
		if err != nil {
			log.Fatal("error reading the response", err)
		}
		prettyPrint(body)
	case "check":
		// check/verify the rules
		fmt.Println("check")
		body, e := CheckCurrentRules()
		if e != nil {
			log.Fatal(e)
		}
		prettyPrint(body)
	case "delete":
		// TODO - implement delete (get the ids from the command line args)
		idsToDelete := []string{
			"1295539185877692419",
			"1295883610038374402",
		}

		body, err := DeleteStreamRules(idsToDelete)
		if err != nil {
			log.Fatal(err)
		}
		prettyPrint(body)
	case "delete-all":
		// TODO - implement delete all
		fmt.Println("delete-all")
	case "help":
		// show the available commands / options
		// TODO - implement
		fmt.Println("help")
		return
	case "stream":
		// subscribe to the feed
		// TODO - implement
		fmt.Println("stream")

		ch := make(chan []byte)

		FetchStream(ch)

		select {
		case result := <-ch:
			prettyPrint(result)
		//case <-time.After(time.Second * 10):
		//	fmt.Println("Server is busy.")
		//	<-ch
		}

		//prettyPrint(body)
	default:
		fmt.Println("the available commands are `add`, `check`, and `stream`")
		os.Exit(1)
	}
}

// FetchStream gets the twitter stream of tweets that match the current rules
func FetchStream(ch chan<- []byte) {
	bearerToken := "Bearer " + apiToken
	req, _ := http.NewRequest(http.MethodGet, streamURL, nil)
	req.Header.Add("Authorization", bearerToken)
	req.Header.Add("Content-type", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)

	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)

	for {
		data, err := read(*reader)

		if err != nil {
			return
		}
		if len(data) == 0 {
			//
			fmt.Println("data is empty")
		}
		prettyPrint(data)

		//select {
		//// send messages, data, or errors
		//case ch <- data:
		//	fmt.Println("sent data")
		//	continue
		////case <-"done":
		////	return
		//}
	}

	//return data, err
	//prettyPrint(data)
}

func read(reader bufio.Reader) ([]byte, error) {
	buffer := new(bytes.Buffer)

	//for {
		line, err := reader.ReadBytes('\n')
		//prettyPrint(line)

		if err != nil && err != io.EOF {
			// all errors other than the end of file error
			return nil, err
		}
		if err == io.EOF && len(line) == 0 {
			if buffer.Len() == 0 {
				return nil, err
			}
			//break
		}
		buffer.Write(line)
	//}
	return buffer.Bytes(), nil
}

func prettyPrint(data []byte) {
	// TODO - clean this up
	var rules bytes.Buffer
	if err := json.Indent(&rules, data,"","\t"); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", string(rules.Bytes()))
}

// CheckCurrentRules fetches the current rules that are persisted
func CheckCurrentRules() ([]byte, error) {
	bearerToken := "Bearer " + apiToken
	req, err := http.NewRequest(http.MethodGet, rulesURL, nil)
	req.Header.Add("Authorization", bearerToken)
	req.Header.Add("Content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error fetching the feed rules: %v", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	return body, err
}

// AddRules adds new rules for the stream. `dryRun` is used to verify the rules, but not persist them
func AddRules(jsonBody []byte, dryRun bool) ([]byte, error) {
	bearerToken := "Bearer " + apiToken

	url := rulesURL
	if dryRun {
		url = url + "?dry_run=true"
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	req.Header.Add("Authorization", bearerToken)
	req.Header.Add("Content-type", "application/json")

	//fmt.Println("req: ", req)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error adding the rules: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	return body, err
}

// DeleteStreamRules deletes persisted rules by rule id
func DeleteStreamRules(ruleIDs DeleteIDs) ([]byte, error) {
	fmt.Println("deleteRules: ", ruleIDs)
	if len(ruleIDs) == 0 {
		return nil, errors.New("you must pass in stream rule ids to delete")
	}

	bearerToken := "Bearer " + apiToken

	ids := map[string]DeleteIDs{"ids": ruleIDs}

	rulesToDelete := DeleteRules{Delete: ids}

	rulesToDeleteJSON, err := json.Marshal(rulesToDelete)
	if err != nil {
		return nil, fmt.Errorf("error converting the rules to a slice of bytes: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, rulesURL, bytes.NewBuffer(rulesToDeleteJSON))
	req.Header.Add("Authorization", bearerToken)
	req.Header.Add("Content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("error deleting the rules: %v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	return body, nil
}
