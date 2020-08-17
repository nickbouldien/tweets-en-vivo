package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const ApiKey = "API_KEY"
const ApiSecret = "API_SECRET"
const ApiToken = "API_TOKEN"

//const twitterApiVersion = 2
const baseURL = "https://api.twitter.com/2"

//var twitterRules []string
var apiToken string

// endpoints
// GET /2/tweets/search/stream/rules
//

type Tweet struct {
	Data []TweetData `json:"data"`
}

type TweetData struct {
	//CreatedAt string `json:"created_at"`
	//FullText string `json:"full_text"`
	ID string `json:"id"`
	Text string `json:"text"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading the .env file")
	}

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

	//updateRules(&twitterRules)

	url := fmt.Sprintf(baseURL + "/tweets?ids=1261326399320715264")
	//url := "http://example.com/"
	fmt.Println("url: ", url)

	body, err := fetch(url)
	if err != nil {
		fmt.Errorf("error reading the response: %v", err)
	}
	fmt.Println("body: ", body)

	var tweet Tweet

	if err := json.Unmarshal([]byte(body), &tweet); err != nil {
		log.Fatal(err)
	}
	fmt.Println("tweet: ", tweet)
}

func fetch(url string) ([]byte, error) {
	bearerToken := "Bearer " + apiToken

	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Authorization", bearerToken)
	//resp, err := http.Get(url)

	fmt.Println("req: ", req)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error fetching the url: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	return body, err
}
