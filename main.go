package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

const ApiKey = "API_KEY"
const ApiSecret = "API_SECRET"
const ApiToken = "API_TOKEN"

const baseURL = "https://api.twitter.com/2"
const rulesURL = baseURL + "/tweets/search/stream/rules"
const streamURL = baseURL + "/tweets/search/stream"

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

	apiKey := os.Getenv(ApiKey)
	apiSecret := os.Getenv(ApiSecret)
	apiToken = os.Getenv(ApiToken)

	// TODO - remove these println statements
	fmt.Println("apiKey: ", apiKey)
	fmt.Println("apiSecret: ", apiSecret)
	fmt.Println("apiToken: ", apiToken)

	if apiKey == "" || apiSecret == "" || apiToken == "" {
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
		prettyPrintRules(body)
	case "check":
		// check/verify the rules
		fmt.Println("check")
		body, e := CheckCurrentRules()
		if e != nil {
			log.Fatal(e)
		}

		prettyPrintRules(body)
	case "stream":
		// subscribe to the feed
		// TODO - implement
		fmt.Println("stream")
		return
	case "delete":
		// TODO - implement delete (get the ids from the command line args)
		idsToDelete := []string{
			"1295539185877692417",
			"1295539185877692418",
		}

		body, e := DeleteStreamRules(idsToDelete)
		if e != nil {
			log.Fatal(e)
		}
		prettyPrintRules(body)
	default:
		fmt.Println("the available commands are `add`, `check`, and `stream`")
		os.Exit(1)
	}
}

func prettyPrintRules(body []byte) {
	// TODO - clean this up
	var rules bytes.Buffer
	if err := json.Indent(&rules, body,"","\t"); err != nil {
		log.Fatal(err)
	}
	//error := json.Indent(&rules, body, "", "\t")

	//prettyJSON, err := json.MarshalIndent(rules, "", "    ")
	//if error != nil {
	//	log.Fatal("Failed to generate json", err)
	//}
	fmt.Printf("%s\n", string(rules.Bytes()))
}

// FetchStream gets the twitter stream of tweets that match the current rules
func FetchStream() ([]byte, error) {
	bearerToken := "Bearer " + apiToken
	req, err := http.NewRequest(http.MethodGet, streamURL, nil)
	req.Header.Add("Authorization", bearerToken)
	req.Header.Add("Content-type", "application/json")

	if err != nil {
		return nil, fmt.Errorf("error fetching the stream: %v", err)
	}

	// TODO - implement
	return nil, errors.New("error: not implemented")
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

	fmt.Println("req: ", req)

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
