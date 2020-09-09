# tweets en vivo

stream tweets to your terminal based on rules (using the [Twitter v2 API](https://developer.twitter.com/en/docs/twitter-api/early-access))

(I know there are good alternatives to this such as [tweetdeck]([https://tweetdeck.twitter.com/]) or for a 
Go library, [go-twitter](https://github.com/dghubble/go-twitter/). This project is just for fun.)


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

(NOTE: all the below examples assume that you have built a binary named `tweets-en-vivo`)

### general usage:
```bash
./tweets-en-vivo -command=<command> [other parameters]
```

### details and examples for each command:

#### add
this command allows you to add rules to the stream.  NOTE: passing a file is optional. the default file is `rules/rules.json`.
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

## project setup
### backend
1. create a [Twitter dev account](https://developer.twitter.com/en/apply-for-access) and project to get an API token 
2. pull down this repo
3. add a `.env` file in the root directory and add your Twitter API token
4. run `go build`

### frontend web client (optional)
check out the [frontend README](./frontend/README.md)


### using docker
```bash
make build
```

```bash
# example running the app in docker with the "stream" command (with the websocket set up on port 5000)
docker run --rm -ti -p 5000:5000 tweets-en-vivo -command=stream -websocket
```

## stream rules
the rules must follow Twitter's documentation. You can put your rules in the `/rules` directory and if desired,
you can place "private rules" (not tracked by git) in the `/rules/private/` directory.


### twitter API documentation
- [filtered streams - rules](https://developer.twitter.com/en/docs/twitter-api/tweets/filtered-stream/integrate/build-a-rule)


### twitter API endpoints
v2 endpoints (base url = https://api.twitter.com/2/)
- `GET /tweets/search/stream`
- `GET /tweets/search/stream/rules`
- `POST /tweets/search/stream/rules`


### resources
- [go-twitter](https://github.com/dghubble/go-twitter)
- [JustForFunc youtube channel](https://www.youtube.com/c/JustForFunc/videos)

### TODOs
There are lots of things I plan on adding/fixing/refactoring.

Here are a few:
- write tests!
- customize the tweet fields retrieved (right now it is decently hard coded to fields I care about)
- customize the tweet filter params (language, location, etc.)
- need better file/code organization
- lots of refactoring/cleanup
- make the websocket server cancelable
- add "hooks" to intercept a tweet and do something with it
- clean up the http requests (there is a lot of duplication)
- make it more easily deployable (use env vars for things like ports, urls, tokens, etc.)
- add the ability to add/delete stream rules from the frontend
- display more tweet fields (tags, etc.) on the frontend

### License

[MIT License](LICENSE)
