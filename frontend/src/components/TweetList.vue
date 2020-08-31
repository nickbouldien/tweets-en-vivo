<template>
  <div id="tweet-list">
    <h2>Tweets</h2>

    <ConnectionInfo
      :error=error
      :websocketOpen=Boolean(socket)
      @toggleWebsocket="toggleConnection"
    />

    <!-- TODO - have a section to display the current stream/tweet "rules" -->

    <ul class="tweets">
      <Tweet
        v-for="tweet in tweets"
        :key="tweet.id"
        :tweet="tweet"
      />
    </ul>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import { Component, Watch } from 'vue-property-decorator';

import ConnectionInfo from '@/components/ConnectionInfo.vue'
import Tweet from '@/components/Tweet.vue'
import { websocketUrl } from '@/config';
import { ITweet, ITweetResponse } from '@/store/main/state';
import { tweetResponses } from '@/data';

@Component({
  components: {
    ConnectionInfo,
    Tweet,
  }
})
export default class TweetList extends Vue {
  error: Error | null = null;
  socket: WebSocket | null = null;
  tweets: ITweet[] = [];

  mounted() {
    // TODO - remove this. this is only for dev/debugging
    this.tweets = tweetResponses.map(tweetResponse => this.mapTweetResponseToTweet(tweetResponse));
  }

  toggleConnection() {
    if (this.socket) {
      // close the socket if it already exists
      console.log("this.socket exists ", this.socket);
      this.socket.close();
      this.socket = null;
    }
    this.createWebSocketConnection();
  }

  mapTweetResponseToTweet(tweetResponse: ITweetResponse): ITweet {
    let author = tweetResponse.includes.users && tweetResponse.includes.users[0];
    let tweetData = tweetResponse.data;

    let tweet: ITweet = {
      id: tweetData.id,
      text: tweetData.text,
      tag: tweetData.tag,
      authorId: tweetData.id,
      createdAt: tweetData.createdAt,
      authorUsername: author.username,
      authorName: author.name,
      matchingRules: tweetResponse.matching_rules,
    }

    return tweet;
  }

  createWebSocketConnection(): void {
    console.log("=> starting the connection to the websocket");
    this.socket = new WebSocket(websocketUrl);

    // register websocket event listeners/handlers
    this.socket.onopen = (event: Event) => {
      console.log('socket opened: ', this.socket);
      this.error = null;
    }

    this.socket.onmessage = (event: MessageEvent) => {
      let tweetResponse: ITweetResponse;

      try {
        tweetResponse = JSON.parse(event.data)
        let tweet: ITweet = this.mapTweetResponseToTweet(tweetResponse);
        console.log("==> tweet: ", tweet);
        // this.tweets.push(tweet);
        // FIXME - better (more efficient) way to do this?
        this.tweets.unshift(tweet);
      } catch(err) {
        console.error("error parsing the websocket message");
        this.error = err;
      }
    }

    this.socket.onerror = (error: Event) => {
      console.error("An error occured with the websocket: ", error);
      this.error = new Error(`An error occured with the websocket`);
    }

    this.socket.onclose = (event: CloseEvent) => {
      console.warn("--> socket closed");
      this.error = null;
      this.socket = null;
      console.warn("--> this.socket: ", this.socket);
    }
  }
}
</script>

<style scoped lang="scss">
#tweet-list {
  padding: 12px;
  text-align: left;
  border: 1px solid black;
}

ul {
  padding-left: 0;
  list-style: none;
}
</style>
