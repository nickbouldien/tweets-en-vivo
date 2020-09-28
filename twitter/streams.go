package twitter

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	wsClient "tweets-en-vivo/websocket"

	"github.com/logrusorgru/aurora"
	"golang.org/x/oauth2"
)

const (
	baseURL    = "https://api.twitter.com/2/tweets/search/stream"
	rulesURL   = baseURL + "/rules"
	streamURL  = baseURL + "?tweet.fields=created_at,lang&expansions=author_id"
	twitterURL = "https://twitter.com/"
)

// the (spoken) languages that you want to see tweets in
// they are BCP47 language tags and are only returned in the Tweet if detected by Twitter
var acceptedLangs = []string{"en", "es", "pt"}

type SimpleTweet struct {
	ID    string `json:"id"`
	Tag   string `json:"tag,omitempty"`
	Value string `json:"value"`
}

type StreamTweet struct {
	AuthorID  string `json:"author_id"`
	CreatedAt string `json:"created_at"`
	Lang      string `json:"lang,omitempty"`
	ID        string `json:"id"`
	Text      string `json:"text"`
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
			fmt.Println("...")
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
	httpClient *http.Client
}

// NewClient creates a new Client
func NewClient(ctx context.Context, token string) *Client {
	return &Client{
		httpClient: oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: token,
			TokenType:   "Bearer",
		})),
	}
}

// FetchStream gets the main stream of tweets that match the current rules
func (client *Client) FetchStream(ch chan<- []byte, done chan<- bool) {
	req, err := http.NewRequest(http.MethodGet, streamURL, nil)
	if err != nil {
		_ = fmt.Errorf("error creating the FetchStream request: %v", err)
	}

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

	done <- true
	return
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

// HandleTweetData handles the tweet data by printing it and sending it to the websocket if the websocket channel is open
func HandleTweetData(wsStream *wsClient.Stream, data []byte) {
	var streamData StreamData
	err := json.Unmarshal(data, &streamData)
	if err != nil {
		_ = fmt.Errorf("error converting the data to a Tweet: %v", err)
		log.Fatal(err)
	}

	// return (skip the tweet) if the tweet isn't in a language you want to see
	if !isAcceptedLanguage(streamData.Data.Lang) {
		return
	}

	// TODO - check if this is safe based on the twitter docs
	author := streamData.Includes["users"][0]

	userURL := fmt.Sprint(twitterURL, author.Username)
	tweetURL := fmt.Sprint(twitterURL, author.Username, "/status/", streamData.Data.ID)

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

	// print the formatted tweet to the terminal
	t.Print()
}

// isAcceptedLanguage returns true if the tweet's language is in `acceptedLangs`
/*
	TODO - check if the twitter API can do this.  I thought it was possible to subscribe to tweets in certain
	languages, but didn't see where that was possible with the v2 API, so doing the filtering here
*/
func isAcceptedLanguage(tweetLanguage string) bool {
	for _, lang := range acceptedLangs {
		if tweetLanguage == lang {
			return true
		}
	}
	return false
}
