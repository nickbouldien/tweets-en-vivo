package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

const ApiToken = "API_TOKEN"

type CheckRulesResponse struct {
	Data []Tweet `json:"data"`
	Meta map[string]string `json:"meta"`
}

type Client struct {
	ApiToken string
	httpClient *http.Client
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(`Error loading the .env file. make sure it exists and 
				has all of the required environment variables`)
	}

	command := flag.String("command", "stream", "the command you want to perform")
	file := flag.String("file", "rules.json", "the rules file you want to use (in the rules/ dir)")
	flag.Parse()

	fmt.Println("--> file:", *file)
	fmt.Println("--> command:", *command)

	//if flag.NArg() < 1 {
	//	fmt.Printf("usage:\n\t%s \"tweets-en-vivo\"\n", os.Args[0])
	//	os.Exit(1)
	//}

	client := Client{
		 os.Getenv(ApiToken),
		 &http.Client{},
	}

	if client.ApiToken == "" {
		log.Fatal(`make sure that you have filled in the required
				main api credentials in the .env file`)
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
		body, err := AddRules(client, byteValue, false)
		if err != nil {
			log.Fatal("error reading the response", err)
		}
		PrettyPrint(body)
	case "check":
		// check/verify the rules
		fmt.Println("check")
		body, e := CheckCurrentRules(client)
		if e != nil {
			log.Fatal(e)
		}
		PrettyPrint(body)
	case "delete":
		// TODO - implement delete (get the ids from the command line args)
		idsToDelete := []string{
			"1295539185877692419",
			"1295883610038374402",
		}

		body, err := DeleteStreamRules(client, idsToDelete)
		if err != nil {
			log.Fatal(err)
		}
		PrettyPrint(body)
	case "delete-all":
		// TODO - implement delete all
		fmt.Println("delete-all")
		// get all the current rule ids
		body, e := CheckCurrentRules(client)
		if e != nil {
			log.Fatal(e)
		}

		var checkResponse CheckRulesResponse

		PrettyPrint(body)

		err := json.Unmarshal(body, &checkResponse)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("checkResponse: ", checkResponse)

		var idsToDelete TweetIDs
		for i, v := range checkResponse.Data {
			fmt.Printf("i: %d, v: %v", i, v)
			idsToDelete = append(idsToDelete, v.ID)
		}
		fmt.Println("idsToDelete: ", idsToDelete)

		resBody, err := DeleteStreamRules(client, idsToDelete)
		if err != nil {
			log.Fatal(err)
		}
		PrettyPrint(resBody)
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

		FetchStream(client, ch)

		select {
		case result := <-ch:
			PrettyPrint(result)
		//case <-"done":
		// TODO - implement
		//	fmt.Println("ending stream.")
		//	close(ch)
		}

		//prettyPrint(body)
	default:
		fmt.Println("--> the available commands are `add`, `check`, `delete`, `delete-all`, and `stream`")
		os.Exit(1)
	}
}
