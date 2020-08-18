package main

import (
	"bytes"
	"encoding/json"
	"errors"
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

//var twitterRules Rules
var apiToken string
var filename string

// just using a slice of bytes for now. the twitter API will validate the rules
//type Rules []byte

type DeleteIDs []string

type DeleteRules struct {
	Delete map[string]DeleteIDs `json:"delete"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(`Error loading the .env file. make sure it exists and 
				has all of the required environment variables`)
	}

	// TODO - get the filename from command line args
	filename = "rules.json"

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

	// 1 - import the `rules.json` file
	jsonFile, err := os.Open(path.Join("rules/", filename))
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	//byteValue, _ := ioutil.ReadAll(jsonFile)
	//fmt.Println("twitterRules: ", twitterRules)

	// 2 - add the rules
	//body, err := AddRules(byteValue, false)
	//if err != nil {
	//	fmt.Errorf("error reading the response: %v", err)
	//}
	//


	// 3 - check/verify the rules

	// 4 - subscribe to the feed

	//idsToDelete := []string{
	//	"1295539185877692417",
	//	"1295539185877692418",
	//}
	//
	//body, e := DeleteStreamRules(idsToDelete)
	//if e != nil {
	//	log.Fatal(e)
	//}

	body, e := CheckCurrentRules()
	if e != nil {
		log.Fatal(e)
	}

	prettyPrintRules(body)
}

func prettyPrintRules(body []byte) {
	var rules interface{}
	if err := json.Unmarshal(body, &rules); err != nil {
		log.Fatal(err)
	}
	fmt.Println("rules: ", rules)
}

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

	// TODO - implement
	//return nil, errors.New("error: not implemented")

	//if err != nil {
	//	return nil, fmt.Errorf("error fetching the rules: %v", err)
	//}
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

	//fmt.Println("rulesToDeleteJSON: ", string(rulesToDeleteJSON))

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
