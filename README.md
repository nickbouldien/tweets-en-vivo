# tweets en vivo

stream tweets to your terminal based on rules

## twitter API documentation
- [filtered streams - rules](https://developer.twitter.com/en/docs/twitter-api/tweets/filtered-stream/integrate/build-a-rule)
- 


## twitter API endpoints
v2 endpoints
- `GET /tweets/search/stream`
- `GET /tweets/search/stream/rules`
- `POST /tweets/search/stream/rules`


## TODOs
- write tests!
- format the printing to terminal
- create the websocket server
- create the frontend
- lots of refactoring/cleanup
- refactor all of the tweet/response typings
- make the websocket server cancelable

## commands
```bash
go build
```

```bash
go mod tidy
```
