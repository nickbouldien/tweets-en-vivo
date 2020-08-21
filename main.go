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
	dryRun := flag.Bool("dryRun", false, "true if you want to verify (but not persist) the rules, otherwise false")

	flag.Parse()
	ruleIDs := flag.Args()

	fmt.Println("--> file:", *file)
	fmt.Println("--> command:", *command)
	fmt.Println("--> ruleIDs:", ruleIDs)
	fmt.Println("--> dryRun:", dryRun)

	token := os.Getenv(ApiToken)
	if token == "" {
		log.Fatal(`make sure that you have filled in the required
				main api credentials in the .env file`)
	}

	client := Client{
		 fmt.Sprint("Bearer ", token),
		 &http.Client{},
	}

	switch *command {
	case "add":
		fmt.Println("add")
		handleAddRulesCommand(client, *file, *dryRun)
	case "check":
		// check/verify the rules
		fmt.Println("check")
		handleCheckRulesCommand(client)
	case "delete":
		// delete the rules with ids passed in as args
		fmt.Println("delete")
		handleDeleteCommand(client, ruleIDs)
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

	go client.FetchStream(ch)

	for  {
		select {
		case data := <-ch:
			handleTweetData(data)
		//case <-"done":
			// TODO - implement
			//	fmt.Println("ending the stream")
			//	close(ch)
		}
	}
}

func handleTweetData(data []byte) {
	// TODO - implement
	PrettyPrint(data)
}

func handleAddRulesCommand(client Client, file string, dryRun bool) {
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
	body, err := client.AddRules(byteValue, dryRun)
	if err != nil {
		log.Fatal("error reading the response", err)
	}
	PrettyPrint(body)
}

func handleCheckRulesCommand(client Client) {
	body, e := client.CheckCurrentRules()
	if e != nil {
		log.Fatal(e)
	}
	PrettyPrint(body)
}

func handleDeleteCommand(client Client, ids TweetIDs) {
	if len(ids) == 0 {
		log.Fatal("you must supply a list of rule ids to delete")
	}

	body, err := client.DeleteStreamRules(ids)
	if err != nil {
		log.Fatal(err)
	}
	PrettyPrint(body)
}

func handleDeleteAllCommand(client Client) {
	// first: get all the current rule ids
	body, e := client.CheckCurrentRules()
	if e != nil {
		log.Fatal(e)
	}

	var currentStreamRules CheckRulesResponse

	PrettyPrint(body)

	err := json.Unmarshal(body, &currentStreamRules)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("currentStreamRules: ", currentStreamRules)

	var idsToDelete TweetIDs
	for i, v := range currentStreamRules.Data {
		fmt.Printf("i: %d, v: %v", i, v)
		idsToDelete = append(idsToDelete, v.ID)
	}
	fmt.Println("idsToDelete: ", idsToDelete)

	resBody, err := client.DeleteStreamRules(idsToDelete)
	if err != nil {
		log.Fatal(err)
	}
	PrettyPrint(resBody)
}

func handleHelpCommand() {
	// TODO - implement (show all of the commands / args)
	fmt.Println("help")
}
