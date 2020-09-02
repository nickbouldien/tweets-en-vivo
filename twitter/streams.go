package twitter

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
