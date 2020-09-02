package CLI

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
	"tweets-en-vivo/twitter"
	"tweets-en-vivo/util"
	wsClient "tweets-en-vivo/websocket"

	"github.com/gorilla/websocket"
)

type Options struct {
	Command         string
	CreateWebsocket bool
	DryRun          bool
	File            string
	RuleIDs         []string
}

func HandleCLICommand(options Options, wg *sync.WaitGroup) {
	apiToken := os.Getenv("API_TOKEN")
	if apiToken == "" {
		log.Fatal(`make sure that you have filled in the required
				api credentials in the .env file`)
	}

	token := fmt.Sprint("Bearer ", apiToken)
	client := twitter.NewClient(token)

	switch options.Command {
	case "add":
		handleAddRulesCommand(client, options.File, options.DryRun)
	case "check":
		// check (fetch) the current rules
		handleCheckRulesCommand(client)
	case "delete":
		// delete the rules with ids passed in as args
		handleDeleteCommand(client, options.RuleIDs)
	case "delete-all":
		// delete all of the current rules
		handleDeleteAllCommand(client)
	case "help":
		// show the available commands / options
		handleHelpCommand()
		return
	case "stream":
		// subscribe to the feed
		fmt.Println("createWebsocket: ", options.CreateWebsocket)

		// FIXME - clean all of this up

		if options.CreateWebsocket {
			wg.Add(1)
			// only start the websocket connection if the -websocket arg is present
			http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
				ws, err := wsClient.Upgrader.Upgrade(w, r, nil)

				if err != nil {
					if _, ok := err.(websocket.HandshakeError); !ok {
						_ = fmt.Errorf("websocket handshake error : %v", err)
					}
					_ = fmt.Errorf("websocket error: %v", err)
					return
				}

				// just for a sanity check
				fmt.Println("websocket addr: ", ws.LocalAddr())

				websocketStream := wsClient.NewStream(ws, make(chan []byte))

				// TODO - fix this
				handleStreamCommand(client, websocketStream)
			})

			go func() {
				log.Fatal(http.ListenAndServe(wsClient.Addr, nil))
			}()
		} else {
			handleStreamCommand(client, nil)
		}
	default:
		handleHelpCommand()
		os.Exit(1)
	}
}

func handleAddRulesCommand(client *twitter.Client, filename string, dryRun bool) {
	// first: import the rules json file
	file, err := os.Open(path.Join("rules/", filename))
	if err != nil {
		log.Fatal("could not open the json file", err)
	}
	defer util.CloseFile(file)

	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("could not read the json", err)
	}

	// second: add the rules
	rules, err := client.AddRules(b, dryRun)
	if err != nil {
		log.Fatal("error reading the response", err)
	}
	util.PrettyPrint(rules)
}

func handleCheckRulesCommand(client *twitter.Client) {
	rules, err := client.FetchCurrentRules()
	if err != nil {
		log.Fatal(err)
	}
	util.PrettyPrint(rules)
}

func handleDeleteCommand(client *twitter.Client, ids twitter.TweetIDs) {
	if len(ids) == 0 {
		log.Fatal("you must supply a list of rule ids to delete")
	}

	rules, err := client.DeleteStreamRules(ids)
	if err != nil {
		log.Fatal(err)
	}
	util.PrettyPrint(rules)
}

func handleDeleteAllCommand(client *twitter.Client) {
	// first: get all the current rule ids
	currentRules, e := client.FetchCurrentRules()
	if e != nil {
		log.Fatal(e)
	}

	printData := map[string][]twitter.SimpleTweet{
		"currentRules": currentRules.Data,
	}

	util.PrettyPrint(printData)

	var idsToDelete twitter.TweetIDs
	for _, v := range currentRules.Data {
		idsToDelete = append(idsToDelete, v.ID)
	}

	// second: use the current rule ids to delete the current rules
	rules, err := client.DeleteStreamRules(idsToDelete)
	if err != nil {
		log.Fatal(err)
	}
	util.PrettyPrint(rules)
}

func handleHelpCommand() {
	fmt.Println(`--> the available commands are: 
				"add", "check", "delete", "delete-all", and "stream"`)
}

func handleStreamCommand(client *twitter.Client, wsStream *wsClient.Stream) {
	if wsStream != nil {
		fmt.Println("there is a websocket connection open")
		go wsStream.Handler(wsStream.WsChannel)
	}

	ch := make(chan []byte)
	go client.FetchStream(ch)

	for {
		select {
		case data := <-ch:
			twitter.HandleTweetData(wsStream, data)
			//case <-"done":
			// TODO - implement
			//	fmt.Println("ending the stream")
			//	close(ch)
			//default:
			//	fmt.Println("default")
		}
	}
}
