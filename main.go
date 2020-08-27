package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

const ApiToken = "API_TOKEN"

//const pingPeriod = 3 * time.Second

type Client struct {
	ApiToken   string
	httpClient *http.Client
	ws         *websocket.Conn
	wsChannel  chan []byte
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// FIXME - range over a list of accepted origins
		return true
	},
}

var addr = ":5000"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(`Error loading the .env file. Make sure it exists and 
				has all of the required environment variables`)
	}
	//interrupt := make(chan os.Signal, 1)

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

	fmt.Println("--> file:", *file)
	fmt.Println("--> command:", *command)

	token := os.Getenv(ApiToken)
	if token == "" {
		log.Fatal(`make sure that you have filled in the required
				api credentials in the .env file`)
	}

	client := Client{
		fmt.Sprint("Bearer ", token),
		&http.Client{
			Timeout: 15 * time.Second,
		},
		nil,
		nil,
	}

	var wg sync.WaitGroup
	wg.Add(1)

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
		fmt.Println("createWebsocket ", *createWebsocket)

		if *createWebsocket {
			// only start the websocket connection if the -websocket arg is present
			http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
				fmt.Println("starting up websocket")
				ws, err := upgrader.Upgrade(w, r, nil)

				if err != nil {
					if _, ok := err.(websocket.HandshakeError); !ok {
						_ = fmt.Errorf("websocket handshake error : %v", err)
					}
					_ = fmt.Errorf("websocket error: %v", err)
					return
				}

				// just for a sanity check
				fmt.Println(ws.LocalAddr())

				// update the client with the websocket and websocket channel
				client.ws = ws
				client.wsChannel = make(chan []byte)

				wg.Add(1)

				handleStreamCommand(client, &wg)
			})

			go func() {
				fmt.Println("http listen and serve")
				wg.Add(1)

				log.Fatal(http.ListenAndServe(addr, nil))
			}()
			//go func() {
			//	handleStreamCommand(client, &wg)
			//}()
		} else {
			handleStreamCommand(client, &wg)
		}
	default:
		fmt.Println(`--> the available commands are: 
				"add", "check", "delete", "delete-all", and "stream"`)
		os.Exit(1)
	}
	wg.Wait()
}

func websocketWriter(ws *websocket.Conn, ch <-chan []byte) {
	defer func() {
		_ = ws.Close()
	}()

	for {
		select {
		case data := <-ch:
			if err := ws.WriteMessage(websocket.TextMessage, data); err != nil {
				_ = fmt.Errorf("error writing the message: %v", err)
				return
			}
		}
	}
}

func handleStreamCommand(client Client, wg *sync.WaitGroup) {
	ch := make(chan []byte, 100)

	wg.Add(1)

	if client.ws != nil && client.wsChannel != nil {
		fmt.Println("there is a websocket connection to send the data to the browser")
		go websocketWriter(client.ws, client.wsChannel)
	}

	go client.FetchStream(ch)

	for {
		select {
		case data := <-ch:
			//return &tweets, nil
			handleTweetData(client, data)
			//case <-"done":
			// TODO - implement
			//	fmt.Println("ending the stream")
			//	close(ch)
		}
	}
}

func handleTweetData(client Client, data []byte) {
	if client.ws != nil && client.wsChannel != nil {
		// if there is an open websocket connection, send the data to it
		client.wsChannel <- data
	}

	//PrettyPrint(data)

	
	var tweet Tweet
	err := json.Unmarshal(data, &tweet)
	if err != nil {
		log.Fatal(err)
	}

	// print to the terminal
	//PrettyPrint(tweet)

	PrettyPrintByteSlice(data)
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
	rules, err := client.AddRules(byteValue, dryRun)
	if err != nil {
		log.Fatal("error reading the response", err)
	}
	PrettyPrint(rules)
}

func handleCheckRulesCommand(client Client) {
	rules, e := client.FetchCurrentRules()
	if e != nil {
		log.Fatal(e)
	}
	PrettyPrint(rules)
	//return rules
}

func handleDeleteCommand(client Client, ids TweetIDs) {
	if len(ids) == 0 {
		log.Fatal("you must supply a list of rule ids to delete")
	}

	rules, err := client.DeleteStreamRules(ids)
	if err != nil {
		log.Fatal(err)
	}
	PrettyPrint(rules)
}

func handleDeleteAllCommand(client Client) {
	// first: get all the current rule ids
	currentRules, e := client.FetchCurrentRules()
	if e != nil {
		log.Fatal(e)
	}

	printData := map[string][]Tweet{
		"currentRules": currentRules.Data,
	}

	PrettyPrint(printData)

	var idsToDelete TweetIDs
	for _, v := range currentRules.Data {
		idsToDelete = append(idsToDelete, v.ID)
	}

	rules, err := client.DeleteStreamRules(idsToDelete)
	if err != nil {
		log.Fatal(err)
	}
	PrettyPrint(rules)
}

func handleHelpCommand() {
	// TODO - implement (show all of the commands / args)
	fmt.Println("help")
}
