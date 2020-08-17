package main

import (
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

//const apiVersion = 2
const baseURL = "https://api.twitter.com/2"

//var twitterRules []string

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading the .env file")
	}

	apiKey := os.Getenv(ApiKey)
	apiSecret := os.Getenv(ApiSecret)
	apiToken := os.Getenv(ApiToken)

	fmt.Println("apiKey: ", apiKey)
	fmt.Println("apiSecret: ", apiSecret)
	fmt.Println("apiToken: ", apiToken)

	//updateRules(&twitterRules)

	url := fmt.Sprintf(baseURL + "/tweets?ids=1261326399320715264")
	fmt.Println("url: ", url)

	//res, err := fetch(url)

}

func fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Errorf("error fetching the url: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	return body, nil
}
