package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

type Client struct {
	ApiToken   string
	httpClient *http.Client
	Ws         *websocket.Conn
	WsChannel  chan []byte
}

const (
	ApiToken      = "API_TOKEN"
	websocketAddr = ":5000"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// FIXME - range over a list of accepted origins
		return true
	},
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

	fmt.Println("--> file:", *file)
	fmt.Println("--> command:", *command)

	token := os.Getenv(ApiToken)
	if token == "" {
		log.Fatal(`make sure that you have filled in the required
				api credentials in the .env file`)
	}

	client := &Client{
		fmt.Sprint("Bearer ", token),
		&http.Client{},
		nil,
		nil,
	}

	var wg sync.WaitGroup

	switch *command {
	case "add":
		handleAddRulesCommand(client, *file, *dryRun)
	case "check":
		// check/verify the rules
		handleCheckRulesCommand(client)
	case "delete":
		// delete the rules with ids passed in as args
		handleDeleteCommand(client, ruleIDs)
	case "delete-all":
		// delete all of the current rules
		handleDeleteAllCommand(client)
	case "help":
		// show the available commands / options
		handleHelpCommand()
		return
	case "stream":
		// subscribe to the feed
		fmt.Println("createWebsocket: ", *createWebsocket)

		if *createWebsocket {
			wg.Add(1)
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
				fmt.Println("websocket addr: ", ws.LocalAddr())

				// update the client struct with the websocket and websocket channel
				client.Ws = ws
				client.WsChannel = make(chan []byte)

				// TODO - fix this
				handleStreamCommand(client)
			})

			go func() {
				log.Fatal(http.ListenAndServe(websocketAddr, nil))
			}()
		} else {
			handleStreamCommand(client)
		}
	default:
		handleHelpCommand()
		os.Exit(1)
	}
	wg.Wait()
}

func websocketWriter(ws *websocket.Conn, ch <-chan []byte) {
	defer func() {
		fmt.Println("closing the websocket connection")
		_ = ws.Close()
	}()

	for {
		select {
		case data := <-ch:
			if err := ws.WriteMessage(websocket.TextMessage, data); err != nil {
				_ = fmt.Errorf("error writing the message to the websocket: %v", err)
				return
			}
		}
	}
}

func handleStreamCommand(client *Client) {
	ch := make(chan []byte)

	if client.Ws != nil && client.WsChannel != nil {
		fmt.Println("there is a websocket connection to send the data to the browser")
		go websocketWriter(client.Ws, client.WsChannel)
	}

	go client.FetchStream(ch)

	for {
		select {
		case data := <-ch:
			handleTweetData(client, data)
			//case <-"done":
			// TODO - implement
			//	fmt.Println("ending the stream")
			//	close(ch)
			//default:
			//	fmt.Println("default")
		}
	}
}

func handleTweetData(client *Client, data []byte) {
	fmt.Println("handleTweetData")
	if client.Ws != nil && client.WsChannel != nil {
		// if there is an open websocket connection, send the data to it
		client.WsChannel <- data
	}

	//var tweet Tweet
	//err := json.Unmarshal(data, &tweet)
	//if err != nil {
	//	log.Fatal(err)
	//}

	PrettyPrintByteSlice(data)
}

func handleAddRulesCommand(client *Client, filename string, dryRun bool) {
	// import the rules json file
	file, err := os.Open(path.Join("rules/", filename))
	if err != nil {
		log.Fatal("could not open the json file", err)
	}
	defer CloseFile(file)

	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("could not read the json", err)
	}

	// add the rules
	rules, err := client.AddRules(b, dryRun)
	if err != nil {
		log.Fatal("error reading the response", err)
	}
	PrettyPrint(rules)
}

func handleCheckRulesCommand(client *Client) {
	rules, err := client.FetchCurrentRules()
	if err != nil {
		log.Fatal(err)
	}
	PrettyPrint(rules)
}

func handleDeleteCommand(client *Client, ids TweetIDs) {
	if len(ids) == 0 {
		log.Fatal("you must supply a list of rule ids to delete")
	}

	rules, err := client.DeleteStreamRules(ids)
	if err != nil {
		log.Fatal(err)
	}
	PrettyPrint(rules)
}

func handleDeleteAllCommand(client *Client) {
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
	fmt.Println(`--> the available commands are: 
				"add", "check", "delete", "delete-all", and "stream"`)
}
