package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"os"
	"path"

	"tweets-en-vivo/twitter"
)

const ApiToken = "API_TOKEN"

var apiToken string

//type DeleteIDs []string

type Tweet struct {
	ID string `json:"id"`
	Value string `json:"value"`
}

type CheckRulesResponse struct {
	Data []Tweet `json:"data"`
	Meta map[string]string `json:"meta"`
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

	apiToken = os.Getenv(ApiToken)

	if apiToken == "" {
		log.Fatal(`make sure that you have filled in the required
				twitter api credentials in the .env file`)
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
		body, err := twitter.AddRules(byteValue, false)
		if err != nil {
			log.Fatal("error reading the response", err)
		}
		twitter.PrettyPrint(body)
	case "check":
		// check/verify the rules
		fmt.Println("check")
		body, e := twitter.CheckCurrentRules()
		if e != nil {
			log.Fatal(e)
		}
		twitter.PrettyPrint(body)
	case "delete":
		// TODO - implement delete (get the ids from the command line args)
		idsToDelete := []string{
			"1295539185877692419",
			"1295883610038374402",
		}

		body, err := twitter.DeleteStreamRules(idsToDelete)
		if err != nil {
			log.Fatal(err)
		}
		twitter.PrettyPrint(body)
	case "delete-all":
		// TODO - implement delete all
		fmt.Println("delete-all")
		// get all the current rule ids
		body, e :=twitter.CheckCurrentRules()
		if e != nil {
			log.Fatal(e)
		}

		var checkResponse CheckRulesResponse

		twitter.PrettyPrint(body)

		err := json.Unmarshal(body, &checkResponse)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("checkResponse: ", checkResponse)

		var idsToDelete twitter.DeleteIDs
		for i, v := range checkResponse.Data {
			fmt.Printf("i: %d, v: %v", i, v)
			idsToDelete = append(idsToDelete, v.ID)
		}
		fmt.Println("idsToDelete: ", idsToDelete)

		resBody, err := twitter.DeleteStreamRules(idsToDelete)
		if err != nil {
			log.Fatal(err)
		}
		twitter.PrettyPrint(resBody)
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

		twitter.FetchStream(ch)

		select {
		case result := <-ch:
			twitter.PrettyPrint(result)
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
