<template>
  <section id="tweet-list">
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
  </section>
</template>

<script lang="ts">
import Vue from 'vue';
import { Component, Watch } from 'vue-property-decorator';

import ConnectionInfo from '@/components/ConnectionInfo.vue'
import Tweet from '@/components/Tweet.vue'
import { websocketUrl } from '@/config';
import { ITweet } from '@/types';
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

  // mounted() {
  //   // TODO - remove this. this is only for dev/debugging
  //   this.tweets = tweetResponses.map(tweetResponse => this.mapTweetResponseToTweet(tweetResponse));
  // }

  toggleConnection() {
    if (this.socket) {
      // close the socket if it already exists
      console.log("this.socket exists: ", this.socket);
      this.socket.close();
      this.socket = null;
    }
    this.createWebSocketConnection();
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
      let tweet: ITweet;

      try {
        tweet = JSON.parse(event.data)
        console.log("==> tweet: ", tweet);
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
  border-radius: 10px;
  padding: 12px;
  text-align: left;
}

ul {
  height: 50vh;
  list-style: none;
  overflow-y: scroll;
  padding-left: 0;
}
</style>
