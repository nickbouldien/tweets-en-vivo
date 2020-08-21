package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	baseURL = "https://api.twitter.com/2"
	rulesURL = streamURL + "/rules"
	streamURL = baseURL + "/tweets/search/stream"
)

// FIXME - refactor all of the structs to make clearer and remove duplication

// TODO - rename this
type CheckRulesResponse struct {
	Data []Tweet `json:"data"`
	Meta map[string]string `json:"meta"`
}

type TweetIDs []string

// TODO - rename this
type DeleteRules struct {
	Delete map[string]TweetIDs `json:"delete"`
}

type Tweet struct {
	ID string `json:"id"`
	Value string `json:"value"`
}

type MatchingRule struct {
	ID string `json:"id"`
	Tag string `json:"tag"`
}

type StreamTweet struct {
	ID string `json:"id"`
	Text string `json:"text"`
}

type StreamData struct {
	Data StreamTweet `json:"data"`
	MatchingRules []MatchingRule `json:"matching_rules"`
}

// FetchStream gets the main stream of tweets that match the current rules
func (client Client) FetchStream(ch chan<- []byte) {
	req, _ := http.NewRequest(http.MethodGet, streamURL, nil)
	req.Header.Add("Authorization", client.ApiToken)

	resp, _ := client.httpClient.Do(req)

	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)

	for {
		data, err := Read(*reader)

		if err != nil {
			_ = fmt.Errorf("error reading the twitter stream: %v", err)
			break
		}

		if len(data) == 0 {
			fmt.Println("the data is empty")
		}

		// send the data over the channel
		ch <- data
	}
}

// CheckCurrentRules fetches the current rules that are persisted
func (client Client) CheckCurrentRules() ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, rulesURL, nil)
	req.Header.Add("Authorization", client.ApiToken)

	resp, err := client.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error fetching the feed rules: %v", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	return body, err
}

// AddRules adds new rules for the stream. `dryRun` is used to verify the rules, but not persist them
func (client Client) AddRules(jsonBody []byte, dryRun bool) ([]byte, error) {
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

	return body, err
}

// DeleteStreamRules deletes persisted rules by rule id
func (client Client) DeleteStreamRules(ruleIDs TweetIDs) ([]byte, error) {
	fmt.Println("deleteRules: ", ruleIDs)
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

	return body, nil
}
