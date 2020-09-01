package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/gorilla/websocket"
)

func handleCLICommand(options Options, wg *sync.WaitGroup) {
	client := &Client{
		fmt.Sprint("Bearer ", ApiToken),
		&http.Client{},
	}

	switch options.command {
	case "add":
		handleAddRulesCommand(client, options.file, options.dryRun)
	case "check":
		// check (fetch) the current rules
		handleCheckRulesCommand(client)
	case "delete":
		// delete the rules with ids passed in as args
		handleDeleteCommand(client, options.ruleIDs)
	case "delete-all":
		// delete all of the current rules
		handleDeleteAllCommand(client)
	case "help":
		// show the available commands / options
		handleHelpCommand()
		return
	case "stream":
		// subscribe to the feed
		fmt.Println("createWebsocket: ", options.createWebsocket)

		if options.createWebsocket {
			wg.Add(1)
			// only start the websocket connection if the -websocket arg is present
			http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
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

				websocketStream := &WebsocketStream{
					Ws:        ws,
					WsChannel: make(chan []byte),
				}

				// TODO - fix this
				handleStreamCommand(client, websocketStream)
			})

			go func() {
				log.Fatal(http.ListenAndServe(websocketAddr, nil))
			}()
		} else {
			handleStreamCommand(client, nil)
		}
	default:
		handleHelpCommand()
		os.Exit(1)
	}
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

	// second: use the current rule ids to delete the current rules
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

func handleStreamCommand(client *Client, wsStream *WebsocketStream) {
	if wsStream != nil {
		fmt.Println("there is a websocket connection open")
		go wsStream.Handler(wsStream.WsChannel)
	}

	ch := make(chan []byte)
	go client.FetchStream(ch)

	for {
		select {
		case data := <-ch:
			handleTweetData(wsStream, data)
			//case <-"done":
			// TODO - implement
			//	fmt.Println("ending the stream")
			//	close(ch)
			//default:
			//	fmt.Println("default")
		}
	}
}
