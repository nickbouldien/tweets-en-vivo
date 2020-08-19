package twitter

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	baseURL = "https://api.twitter.com/2"
	rulesURL = baseURL + "/tweets/search/stream/rules"
	streamURL = baseURL + "/tweets/search/stream"
)

var apiToken string

type DeleteIDs []string

type DeleteRules struct {
	Delete map[string]DeleteIDs `json:"delete"`
}

type Tweet struct {
	ID string `json:"id"`
	Value string `json:"value"`
}

type CheckRulesResponse struct {
	Data []Tweet `json:"data"`
	Meta map[string]string `json:"meta"`
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
		PrettyPrint(data)

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

func PrettyPrint(data []byte) {
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
