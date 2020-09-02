package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
	"tweets-en-vivo/CLI"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(`Error loading the .env file. Make sure it exists and 
				has all of the required environment variables`)
	}

	command := flag.String("command", "stream", "the command you want to perform")
	file := flag.String("file", "rules.json", "the rules file you want to use (in the rules/ dir)")
	dryRun := flag.Bool("dryRun", false, "true if you want to verify (but not persist) the rules, otherwise false")
	createWebsocket := flag.Bool(
		"websocket",
		false,
		"true if you want to create a websocket to send the data to a frontend client, otherwise false",
	)

	flag.Parse()
	ruleIDs := flag.Args()

	options := CLI.Options{
		Command:         *command,
		CreateWebsocket: *createWebsocket,
		File:            *file,
		DryRun:          *dryRun,
		RuleIDs:         ruleIDs,
	}

	fmt.Println("--> file:", options.File)
	fmt.Println("--> command:", options.Command)

	var wg sync.WaitGroup

	CLI.HandleCLICommand(options, &wg)

	wg.Wait()
}
