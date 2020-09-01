package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

const (
	websocketAddr = ":5000"
)

var ApiToken string

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// FIXME - range over a list of accepted origins
		return true
	},
}

type Options struct {
	command         string
	createWebsocket bool
	dryRun          bool
	file            string
	ruleIDs         []string
}

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

	options := Options{
		command:         *command,
		createWebsocket: *createWebsocket,
		file:            *file,
		dryRun:          *dryRun,
		ruleIDs:         ruleIDs,
	}

	fmt.Println("--> file:", options.file)
	fmt.Println("--> command:", options.command)

	ApiToken = os.Getenv("API_TOKEN")
	if ApiToken == "" {
		log.Fatal(`make sure that you have filled in the required
				api credentials in the .env file`)
	}

	var wg sync.WaitGroup

	handleCLICommand(options, &wg)

	wg.Wait()
}
