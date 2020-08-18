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

var twitterRules Rules
var apiToken string
var filename string

type Rules []byte

//type Tweet struct {
//	Data []TweetData `json:"data"`
//}

//type TweetData struct {
//	//CreatedAt string `json:"created_at"`
//	//FullText string `json:"full_text"`
//	ID string `json:"id"`
//	Text string `json:"text"`
//}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading the .env file")
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
	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &twitterRules)

	fmt.Println("twitterRules: ", twitterRules)

	// 2 - add the rules
	body, err := addRules(twitterRules)
	if err != nil {
		fmt.Errorf("error reading the response: %v", err)
	}

	var rules interface{}
	if err := json.Unmarshal(body, &rules); err != nil {
		log.Fatal(err)
	}
	fmt.Println("rules: ", rules)

	// 3 - check/verify the rules

	// 4 - subscribe to the feed
}

func fetchStream() ([]byte, error) {
	bearerToken := "Bearer " + apiToken
	req, err := http.NewRequest(http.MethodGet, streamURL, nil)
	req.Header.Add("Authorization", bearerToken)

	if err != nil {
		return nil, fmt.Errorf("error fetching the stream: %v", err)
	}

	// TODO - implement
	return nil, errors.New("error: not implemented")
}

func checkRules() ([]byte, error) {
	bearerToken := "Bearer " + apiToken
	req, err := http.NewRequest(http.MethodGet, rulesURL, nil)
	req.Header.Add("Authorization", bearerToken)

	if err != nil {
		return nil, fmt.Errorf("error fetching the feed rules: %v", err)
	}

	// TODO - implement
	return nil, errors.New("error: not implemented")
}

func addRules(jsonBody []byte) ([]byte, error) {
	bearerToken := "Bearer " + apiToken

	req, err := http.NewRequest(http.MethodPost, rulesURL, bytes.NewBuffer(jsonBody))
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

func deleteRule(ruleID string) error {
	// TODO - implement
	return errors.New("error: not implemented")
}
