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
		handleAddRulesCommand(client, *file)
	case "check":
		// check/verify the rules
		fmt.Println("check")
		handleCheckRulesCommand(client)
	case "delete":
		// delete the rules with ids passed in as args
		fmt.Println("delete")
		handleDeleteCommand(client)
	case "delete-all":
		// delete all of the current rules
		fmt.Println("delete-all")
		handleDeleteAllCommand(client)
	case "help":
		// show the available commands / options
		handleHelpCommand()
		return
	case "stream":
		// subscribe to the feed
		fmt.Println("stream")
		handleStreamCommand(client)
	default:
		fmt.Println(`--> the available commands are: 
				"add", "check", "delete", "delete-all", and "stream"`)
		os.Exit(1)
	}
}

func handleStreamCommand(client Client) {
	ch := make(chan []byte, 100)

	go FetchStream(client, ch)

	for  {
		select {
		case result := <-ch:
			fmt.Println("got data!!!")
			PrettyPrint(result)
			//case <-"done":
			// TODO - implement
			//	fmt.Println("ending stream.")
			//	close(ch)
		}
	}


	//prettyPrint(body)
}

func handleAddRulesCommand(client Client, file string) {
	// import the rules json file
	jsonFile, err := os.Open(path.Join("rules/", file))
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
}

func handleCheckRulesCommand(client Client) {
	body, e := CheckCurrentRules(client)
	if e != nil {
		log.Fatal(e)
	}
	PrettyPrint(body)
}

func handleDeleteCommand(client Client) {
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
}

func handleDeleteAllCommand(client Client) {
	// first: get all the current rule ids
	body, e := CheckCurrentRules(client)
	if e != nil {
		log.Fatal(e)
	}

	var currentStreamRules CheckRulesResponse

	PrettyPrint(body)

	err := json.Unmarshal(body, &currentStreamRules)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("checkResponse: ", currentStreamRules)

	var idsToDelete TweetIDs
	for i, v := range currentStreamRules.Data {
		fmt.Printf("i: %d, v: %v", i, v)
		idsToDelete = append(idsToDelete, v.ID)
	}
	fmt.Println("idsToDelete: ", idsToDelete)

	resBody, err := DeleteStreamRules(client, idsToDelete)
	if err != nil {
		log.Fatal(err)
	}
	PrettyPrint(resBody)
}

func handleHelpCommand() {
	// TODO - implement
	fmt.Println("help")
}
