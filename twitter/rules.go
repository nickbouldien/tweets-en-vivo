package twitter

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type AddRulesResponse struct {
	Meta AddRulesMeta  `json:"meta"`
	Data []SimpleTweet `json:"data"`
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
	Data []SimpleTweet     `json:"data"`
	Meta map[string]string `json:"meta"`
}

type MatchingRule struct {
	ID  int64  `json:"id"`
	Tag string `json:"tag"`
}

// FetchCurrentRules fetches the current rules that are persisted
func (client *Client) FetchCurrentRules() (*FetchRulesResponse, error) {
	req, err := http.NewRequest(http.MethodGet, rulesURL, nil)
	req.Header.Add("Authorization", client.apiToken)

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
	req.Header.Add("Authorization", client.apiToken)
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
	req.Header.Add("Authorization", client.apiToken)
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
