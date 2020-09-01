# tweets en vivo

stream tweets to your terminal based on rules (using the Twitter v2 API)

> the rules must follow Twitter's documentation. You can put your rules in the `/rules` directory and if desired,
you can place "private rules" (not tracked by git) in the `/rules/private/` directory.

## quick demo
![tweets-en-vivo demo](demo.gif)


## commands
list of available commands:
- add
- check
- stream
- delete <list of rule ids>
- delete-all
- help

(NOTE: all of this is assuming you have built a binary named `tweets-en-vivo`)

### general usage:
```bash
./tweets-en-vivo -command=<command> [other parameters]
```

### details and examples for each command:

#### add
this command allows you to add rules to the stream.  NOTE: passing a file is optional. the default file is `rules.json`
```bash
# assuming the files lives in the `rules/` directory
./tweets-en-vivo -command=add -file=my-rules.json
```

```bash
# this example uses a "private" rules file in the `rules/private/` directory
./tweets-en-vivo -command=add -file=private/investing.json
```

#### check
this command allows you to check the current rules of the stream
```bash
./tweets-en-vivo -command=check
```

#### stream
this command allows you to access the tweet stream
(and optionally create a websocket server to send the tweets to a frontend client)
```bash
# this example shows how to access the stream
./tweets-en-vivo -command=stream
```

```bash
# this example shows how to access the stream and create a websocket server  
./tweets-en-vivo -command=stream -websocket
```

#### delete
this command allows you to delete rules of the stream given rule ids
```bash
./tweets-en-vivo -command=delete "1300496243039318017" "13004962430393180234"
```

#### delete-all
this command allows you to delete all the current rules of the stream
```bash
./tweets-en-vivo -command=delete-all
```

#### help
this command allows you to print out the help "menu" (this simply displays the available commands in the terminal)
NOTE: the help "menu" is also displayed if you type in a command that does not exist
```bash
./tweets-en-vivo -command=help
```


## frontend documentation
[frontend readme](./frontend/README.md)


## twitter API documentation
- [filtered streams - rules](https://developer.twitter.com/en/docs/twitter-api/tweets/filtered-stream/integrate/build-a-rule)
- 

### twitter API endpoints
v2 endpoints (base url = https://api.twitter.com/2/)
- `GET /tweets/search/stream`
- `GET /tweets/search/stream/rules`
- `POST /tweets/search/stream/rules`


## helpful development commands
```bash
go build
```

```bash
go mod tidy
```

## TODOs
- write tests!
- format the printing to terminal
- lots of refactoring/cleanup
- refactor all of the tweet/response typings
- make the websocket server cancelable
