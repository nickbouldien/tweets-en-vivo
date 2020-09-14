package CLI

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"tweets-en-vivo/twitter"
	"tweets-en-vivo/util"
	wsClient "tweets-en-vivo/websocket"

	"github.com/gorilla/websocket"
)

// Options holds all the options the user entered from the command line
type Options struct {
	command         string
	createWebsocket bool
	dryRun          bool
	file            string
	ruleIDs         []string
}

// NewOptions creates a new Options struct
func NewOptions(cmd string, createWebsocket bool, dryRun bool, file string, ruleIds []string) *Options {
	return &Options{
		command:         cmd,
		createWebsocket: createWebsocket,
		dryRun:          dryRun,
		file:            file,
		ruleIDs:         ruleIds,
	}
}

func (o *Options) HandleCommand(wg *sync.WaitGroup) {
	apiToken := os.Getenv("API_TOKEN")
	if apiToken == "" {
		log.Fatal(`make sure that you have filled in the required
				api credentials in the .env file`)
	}

	fmt.Println("--> file:", o.file)
	fmt.Println("--> command:", o.command)

	token := fmt.Sprint("Bearer ", apiToken)
	client := twitter.NewClient(token)

	switch o.command {
	case "add":
		o.handleAddRulesCommand(client)
	case "check":
		// check (fetch) the current rules
		o.handleCheckRulesCommand(client)
	case "delete":
		// delete the rules with ids passed in as args
		o.handleDeleteCommand(client)
	case "delete-all":
		// delete all of the current rules
		o.handleDeleteAllCommand(client)
	case "help":
		// show the available commands / options
		o.handleHelpCommand()
		return
	case "stream":
		// subscribe to the feed
		fmt.Println("createWebsocket: ", o.createWebsocket)

		// FIXME - clean all of this up. find better way to asynchronously run the websocket server and
		// handle the connection with the twitter API

		if o.createWebsocket {
			wg.Add(1)

			portEnvVar := os.Getenv("WEBSOCKET_PORT")
			wsPort := fmt.Sprintf(":%s", portEnvVar)

			if wsPort == ":" {
				// the WEBSOCKET_PORT was empty so use a fallback port
				wsPort = ":5000"
			}

			allowedOriginsEnvVar := os.Getenv("ALLOWED_ORIGINS")
			wsClient.AllowedOrigins = strings.Split(allowedOriginsEnvVar, ",")
			fmt.Println("websocket server allowed origins: ", wsClient.AllowedOrigins)

			// only start the websocket server if the `-websocket` arg is present
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

				o.handleStreamCommand(client, websocketStream)
			})

			go func() {
				log.Fatal(http.ListenAndServe(wsPort, nil))
			}()
		} else {
			o.handleStreamCommand(client, nil)
		}
	default:
		o.handleHelpCommand()
		os.Exit(1)
	}
}

func (o *Options) handleAddRulesCommand(client *twitter.Client) {
	// first: import the rules json file
	file, err := os.Open(path.Join("rules/", o.file))
	if err != nil {
		log.Fatal("could not open the json file", err)
	}
	defer util.CloseFile(file)

	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("could not read the file", err)
	}

	// second: add the rules
	rules, err := client.AddRules(b, o.dryRun)
	if err != nil {
		log.Fatal("error reading the response", err)
	}
	util.PrettyPrint(rules)
}

func (o *Options) handleCheckRulesCommand(client *twitter.Client) {
	rules, err := client.FetchCurrentRules()
	if err != nil {
		log.Fatal(err)
	}
	util.PrettyPrint(rules)
}

func (o *Options) handleDeleteCommand(client *twitter.Client) {
	if len(o.ruleIDs) == 0 {
		log.Fatal("you must supply a list of rule ids to delete")
	}

	rules, err := client.DeleteStreamRules(o.ruleIDs)
	if err != nil {
		log.Fatal(err)
	}
	util.PrettyPrint(rules)
}

func (o *Options) handleDeleteAllCommand(client *twitter.Client) {
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

func (o *Options) handleHelpCommand() {
	fmt.Println(`--> the available commands are: 
				"add", "check", "delete", "delete-all", and "stream"`)
}

func (o *Options) handleStreamCommand(client *twitter.Client, wsStream *wsClient.Stream) {
	if wsStream != nil {
		// start the websocket
		go wsStream.Handler()
	}

	ch := make(chan []byte)
	done := make(chan bool)
	go client.FetchStream(ch, done)

	for {
		select {
		case data := <-ch:
			twitter.HandleTweetData(wsStream, data)
		case <-done:
			// close the channels and the websocket connection
			fmt.Println("ending the stream")
			close(ch)

			wsStream.WSDone <- true
		}
	}
}
