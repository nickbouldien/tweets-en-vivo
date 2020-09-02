package twitter

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	wsClient "tweets-en-vivo/websocket"

	"github.com/logrusorgru/aurora"
)

const (
	baseURL  = "https://api.twitter.com/2/tweets/search/stream"
	rulesURL = baseURL + "/rules"
	// TODO - make the "expanded fields" in the streamURL optional/expandable
	streamURL = baseURL + "?tweet.fields=created_at&expansions=author_id"
)

// FIXME - refactor all of the structs to make clearer and remove duplication

type AddRulesResponse struct {
	Meta AddRulesMeta  `json:"meta"`
	Data []SimpleTweet `json:"data"`
}

type AddRulesSummary struct {
	Created    int64 `json:"created"`
	NotCreated int64 `json:"not_created"`
}

type AddRulesMeta struct {
	Sent    string          `json:"sent"`
	Summary AddRulesSummary `json:"summary"`
}

type DeleteRulesResponse struct {
	Meta DeleteRulesMeta `json:"meta"`
}

type DeleteRulesMeta struct {
	Sent    string             `json:"sent"`
	Summary DeleteRulesSummary `json:"summary"`
}

type DeleteRulesSummary struct {
	Deleted    int64 `json:"deleted"`
	NotDeleted int64 `json:"not_deleted"`
}

type DeleteRules struct {
	Delete map[string]TweetIDs `json:"delete"`
}

type FetchRulesResponse struct {
	Data []SimpleTweet     `json:"data"`
	Meta map[string]string `json:"meta"`
}

type MatchingRule struct {
	ID  int64  `json:"id"`
	Tag string `json:"tag"`
}

type SimpleTweet struct {
	ID    string `json:"id"`
	Tag   string `json:"tag,omitempty"`
	Value string `json:"value"`
}

type StreamTweet struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
	AuthorID  string `json:"author_id"`
}

type StreamData struct {
	Data          StreamTweet       `json:"data"`
	MatchingRules []MatchingRule    `json:"matching_rules"`
	Includes      map[string][]User `json:"includes"`
}

type TweetIDs []string

type Tweet struct {
	AuthorID       string         `json:"authorId"`
	AuthorName     string         `json:"authorName"`
	AuthorUsername string         `json:"authorUsername"`
	CreatedAt      string         `json:"created_at"`
	ID             string         `json:"id"`
	MatchingRules  []MatchingRule `json:"matching_rules"`
	Text           string         `json:"text"`
	TweetURL       string         `json:"tweetUrl"`
	UserURL        string         `json:"userUrl"`
}

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type StreamResponseBodyReader struct {
	reader *bufio.Reader
	buf    bytes.Buffer
}

// Read is a helper function to read the data from the stream
func (r *StreamResponseBodyReader) Read() ([]byte, error) {
	r.buf.Truncate(0)

	for {
		line, err := r.reader.ReadBytes('\n')

		if len(line) == 0 {
			fmt.Println("len(line) == 0")
			continue
		}

		if err != nil && err != io.EOF {
			// all errors other than the end of file error
			_ = fmt.Errorf("read error: %v", err)
			return nil, err
		}

		if err == io.EOF && len(line) == 0 {
			_ = fmt.Errorf("io.EOF && len(line): %v", err)

			if r.buf.Len() == 0 {
				_ = fmt.Errorf("buf.Len() : %v", err)
				return nil, err
			}
			fmt.Println("breaking")
			break
		}

		if bytes.HasSuffix(line, []byte("\r\n")) {
			r.buf.Write(bytes.TrimRight(line, "\r\n"))
			break
		}

		r.buf.Write(line)
	}

	return r.buf.Bytes(), nil
}

// Client connects with the twitter API
type Client struct {
	apiToken   string
	httpClient *http.Client
}

// NewClient creates a new Client
func NewClient(token string) *Client {
	return &Client{
		apiToken:   token,
		httpClient: &http.Client{},
	}
}

// FetchStream gets the main stream of tweets that match the current rules
func (client *Client) FetchStream(ch chan<- []byte) {
	req, err := http.NewRequest(http.MethodGet, streamURL, nil)
	if err != nil {
		_ = fmt.Errorf("error creating the FetchStream request: %v", err)
	}
	req.Header.Add("Authorization", client.apiToken)

	resp, err := client.httpClient.Do(req)
	if err != nil {
		_ = fmt.Errorf("error with the FetchStream request: %v", err)
	}
	defer resp.Body.Close()

	r := StreamResponseBodyReader{reader: bufio.NewReader(resp.Body)}

	for {
		data, err := r.Read()

		if err != nil {
			_ = fmt.Errorf("error reading the twitter stream: %v", err)
			break
		}

		if len(data) == 0 {
			continue
		}

		// send the data over the channel
		ch <- data
	}
}

// FetchCurrentRules fetches the current rules that are persisted
func (client *Client) FetchCurrentRules() (*FetchRulesResponse, error) {
	req, err := http.NewRequest(http.MethodGet, rulesURL, nil)
	req.Header.Add("Authorization", client.apiToken)

	resp, err := client.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error fetching the feed rules: %v", err)
	}

	defer resp.Body.Close()
	fmt.Println("response status code: ", resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)

	var fetchRulesResponse FetchRulesResponse
	err = json.Unmarshal(body, &fetchRulesResponse)
	if err != nil {
		log.Printf("error decoding the response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}
		log.Fatal(err)
	}

	return &fetchRulesResponse, nil
}

// AddRules adds new rules for the stream. `dryRun` is used to verify the rules, but not persist them
func (client *Client) AddRules(jsonBody []byte, dryRun bool) (*AddRulesResponse, error) {
	url := rulesURL
	if dryRun {
		url = url + "?dry_run=true"
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	req.Header.Add("Authorization", client.apiToken)
	req.Header.Add("Content-type", "application/json")

	resp, err := client.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error adding the rules: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var addRulesResponse AddRulesResponse
	err = json.Unmarshal(body, &addRulesResponse)
	if err != nil {
		log.Printf("error decoding the response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}
		log.Fatal(err)
	}

	return &addRulesResponse, err
}

// DeleteStreamRules deletes persisted rules by rule id
func (client *Client) DeleteStreamRules(ruleIDs TweetIDs) (*DeleteRulesResponse, error) {
	fmt.Println("ids to delete: ", ruleIDs)
	if len(ruleIDs) == 0 {
		return nil, errors.New("you must pass in stream rule ids to delete")
	}

	ids := map[string]TweetIDs{"ids": ruleIDs}

	rulesToDelete := DeleteRules{Delete: ids}

	rulesToDeleteJSON, err := json.Marshal(rulesToDelete)
	if err != nil {
		return nil, fmt.Errorf("error converting the rules to a slice of bytes: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, rulesURL, bytes.NewBuffer(rulesToDeleteJSON))
	req.Header.Add("Authorization", client.apiToken)
	req.Header.Add("Content-type", "application/json")

	resp, err := client.httpClient.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("error deleting the rules: %v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	var deleteRulesResponse DeleteRulesResponse
	err = json.Unmarshal(body, &deleteRulesResponse)
	if err != nil {
		log.Printf("error decoding the response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}
		log.Fatal(err)
	}

	return &deleteRulesResponse, nil
}

// Print prints a tweet with formatting/colors
func (t *Tweet) Print() {
	authorNameRunes := []rune(t.AuthorName)
	authorUsernameRunes := []rune(t.AuthorUsername)
	textRunes := []rune(t.Text)

	fmt.Printf("%s - %s\n %s\n %s\n\n",
		aurora.Blue("@"+string(authorUsernameRunes)),
		aurora.Cyan(string(authorNameRunes)),
		aurora.White(string(textRunes)),
		aurora.Underline(aurora.Green(t.TweetURL)),
	)
}

func HandleTweetData(wsStream *wsClient.Stream, data []byte) {
	var streamData StreamData
	err := json.Unmarshal(data, &streamData)
	if err != nil {
		_ = fmt.Errorf("error converting the data to a Tweet: %v", err)
		log.Fatal(err)
	}

	// TODO - check if this is safe based on the twitter docs
	author := streamData.Includes["users"][0]

	userURL := fmt.Sprint("https://twitter.com/", author.Username)
	tweetURL := fmt.Sprint("https://twitter.com/", author.Username, "/status/", streamData.Data.ID)

	t := Tweet{
		AuthorID:       author.ID,
		AuthorName:     author.Name,
		AuthorUsername: author.Username,
		CreatedAt:      streamData.Data.CreatedAt,
		ID:             streamData.Data.ID,
		MatchingRules:  streamData.MatchingRules,
		Text:           streamData.Data.Text,
		TweetURL:       tweetURL,
		UserURL:        userURL,
	}

	b, err := json.Marshal(t)
	if err != nil {
		_ = fmt.Errorf("error marshalling the data to a slice of bytes")
		return
	}

	if wsStream != nil {
		// if there is an open websocket connection, send the data to its channel
		wsStream.WsChannel <- b
	}

	t.Print()
}
