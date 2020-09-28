package main

import (
	"context"
	"flag"
	"log"
	"os"
	"sync"
	"tweets-en-vivo/CLI"
	"tweets-en-vivo/twitter"

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

	apiToken := os.Getenv("API_TOKEN")
	if apiToken == "" {
		log.Fatal(`make sure that you have filled in the required
				api credentials in the .env file`)
	}
	//token := fmt.Sprint("Bearer ", apiToken)

	ctx := context.Background()

	client := twitter.NewClient(ctx, apiToken)

	var wg sync.WaitGroup
	CLIOptions := CLI.NewOptions(*command, *createWebsocket, *dryRun, *file, ruleIDs)
	CLIOptions.HandleCommand(client, &wg)

	wg.Wait()
}
