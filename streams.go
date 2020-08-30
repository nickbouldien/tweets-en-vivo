package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	baseURL  = "https://api.twitter.com/2/tweets/search/stream"
	rulesURL = baseURL + "/rules"
	// TODO - make the "expanded fields" in the streamURL optional/expandable
	streamURL = baseURL + "?tweet.fields=created_at&expansions=author_id"
)

// FIXME - refactor all of the structs to make clearer and remove duplication

type AddRulesResponse struct {
	Meta AddRulesMeta `json:"meta"`
	Data []Tweet      `json:"data"`
}

type AddRulesSummary struct {
	Created    int64 `json:"created"`
	NotCreated int64 `json:"not_created"`
}

type AddRulesMeta struct {
	Sent    string          `json:"sent"`
	Summary AddRulesSummary `json:"summary"`
}

type DeleteRulesResponse struct {
	Meta DeleteRulesMeta `json:"meta"`
}

type DeleteRulesMeta struct {
	Sent    string             `json:"sent"`
	Summary DeleteRulesSummary `json:"summary"`
}

type DeleteRulesSummary struct {
	Deleted    int64 `json:"deleted"`
	NotDeleted int64 `json:"not_deleted"`
}

type DeleteRules struct {
	Delete map[string]TweetIDs `json:"delete"`
}

type FetchRulesResponse struct {
	Data []Tweet           `json:"data"`
	Meta map[string]string `json:"meta"`
}

type TweetIDs []string

type Tweet struct {
	ID    string `json:"id"`
	Tag   string `json:"tag,omitempty"`
	Value string `json:"value"`
}

//type Rule struct {
//	ID    string `json:"id"`
//	Value string `json:"value"`
//}

//type MatchingRule struct {
//	ID  string `json:"id"`
//	Tag string `json:"tag"`
//}

//type StreamTweet struct {
//	ID   string `json:"id"`
//	Text string `json:"text"`
//}

//type StreamData struct {
//	Data          StreamTweet `json:"data"`
//	MatchingRules []Rule      `json:"matching_rules"`
//}

type StreamResponseBodyReader struct {
	reader *bufio.Reader
	buf    bytes.Buffer
}

// FetchStream gets the main stream of tweets that match the current rules
func (client *Client) FetchStream(ch chan<- []byte) {
	req, err := http.NewRequest(http.MethodGet, streamURL, nil)
	if err != nil {
		_ = fmt.Errorf("error creating the FetchStream request: %v", err)
	}
	req.Header.Add("Authorization", client.ApiToken)

	resp, err := client.httpClient.Do(req)
	if err != nil {
		_ = fmt.Errorf("error with the FetchStream request: %v", err)
	}
	defer resp.Body.Close()

	r := StreamResponseBodyReader{reader: bufio.NewReader(resp.Body)}

	for {
		data, err := r.Read()

		if err != nil {
			_ = fmt.Errorf("error reading the twitter stream: %v", err)
			break
		}

		if len(data) == 0 {
			fmt.Println("the data is empty")
			continue
		}

		// send the data over the channel
		ch <- data
	}
}

// CheckCurrentRules fetches the current rules that are persisted
func (client *Client) FetchCurrentRules() (*FetchRulesResponse, error) {
	req, err := http.NewRequest(http.MethodGet, rulesURL, nil)
	req.Header.Add("Authorization", client.ApiToken)

	resp, err := client.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error fetching the feed rules: %v", err)
	}

	defer resp.Body.Close()
	fmt.Println("response status code: ", resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)

	var fetchRulesResponse FetchRulesResponse
	err = json.Unmarshal(body, &fetchRulesResponse)
	if err != nil {
		log.Printf("error decoding the response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}
		log.Fatal(err)
	}

	return &fetchRulesResponse, nil
}

// AddRules adds new rules for the stream. `dryRun` is used to verify the rules, but not persist them
func (client *Client) AddRules(jsonBody []byte, dryRun bool) (*AddRulesResponse, error) {
	url := rulesURL
	if dryRun {
		url = url + "?dry_run=true"
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	req.Header.Add("Authorization", client.ApiToken)
	req.Header.Add("Content-type", "application/json")

	resp, err := client.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error adding the rules: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var addRulesResponse AddRulesResponse
	err = json.Unmarshal(body, &addRulesResponse)
	if err != nil {
		log.Printf("error decoding the response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}
		log.Fatal(err)
	}

	return &addRulesResponse, err
}

// DeleteStreamRules deletes persisted rules by rule id
func (client *Client) DeleteStreamRules(ruleIDs TweetIDs) (*DeleteRulesResponse, error) {
	fmt.Println("ids to delete: ", ruleIDs)
	if len(ruleIDs) == 0 {
		return nil, errors.New("you must pass in stream rule ids to delete")
	}

	ids := map[string]TweetIDs{"ids": ruleIDs}

	rulesToDelete := DeleteRules{Delete: ids}

	rulesToDeleteJSON, err := json.Marshal(rulesToDelete)
	if err != nil {
		return nil, fmt.Errorf("error converting the rules to a slice of bytes: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, rulesURL, bytes.NewBuffer(rulesToDeleteJSON))
	req.Header.Add("Authorization", client.ApiToken)
	req.Header.Add("Content-type", "application/json")

	resp, err := client.httpClient.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("error deleting the rules: %v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	var deleteRulesResponse DeleteRulesResponse
	err = json.Unmarshal(body, &deleteRulesResponse)
	if err != nil {
		log.Printf("error decoding the response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}
		log.Fatal(err)
	}

	return &deleteRulesResponse, nil
}
