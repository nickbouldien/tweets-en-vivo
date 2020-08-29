<template>
  <div id="tweet-list">
    <h2>Tweets</h2>
    <p>
      websocket connection: {{ connectionOpen }}
    </p>
    <div v-if="error != null">
      <pre>
        {{ error.toString() }}
      </pre>
    </div>

    <!-- <button v-on:click="openConnection">Open the websocket connection</button> -->

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

import Tweet from '@/components/Tweet.vue'
import { websocketUrl } from '@/config';
import { ITweet, ITweetResponse } from '@/store/main/state';
import { tweetResponses } from '../data'

@Component({
  components: {
    Tweet,
  }
})
export default class TweetList extends Vue {
  error: Error | null = null;
  tweets: ITweet[] = [];
  socket: WebSocket | null = null;

  mounted() {
    this.error = null;
    this.tweets = tweetResponses.map(tweetResponse => this.mapTweetResponseToTweet(tweetResponse))
    this.socket = null;
  }

  created() {
    // this.createWebSocketConnection();
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
    console.log("Starting connection to WebSocket Server");
    this.socket = new WebSocket(websocketUrl);

    // register websocket event listeners/handlers
    this.socket.onopen = (event: Event) => {
      console.log('socket opened: ', this.socket);
    }

    this.socket.onmessage = (event: MessageEvent) => {
      let tweetResponse: ITweetResponse;
      let tweet: ITweet;
      try {
        console.log("event.data: ", event.data);
        tweetResponse = JSON.parse(event.data)
        
        console.log("tweetResponse: ", tweetResponse);
        tweet = tweetResponse.data;
        let author = tweetResponse.includes.users && tweetResponse.includes.users[0];
        if (author) {
          tweet.authorUsername = author.username;
          tweet.authorName = author.name;
        }
        console.log("tweet: ", tweet);

        this.tweets.push(tweet);
      } catch(err) {
        this.error = err;
      }
    }

    this.socket.onerror = (error: Event) => {
      console.error("An error occured with the websocket: ", error);
      this.error = new Error(`An error occured with the websocket`);
    }

    this.socket.onclose = (event: CloseEvent) => {
      console.warn("--> socket closed");
      this.socket = null;
      console.warn("--> this.socket: ", this.socket);
    }
  }

  get connectionOpen(): string {
    console.log("connectionOpen ", this.socket);
    return this.socket === null ? "open" : "closed";
  }
}
</script>

<style lang="scss">
#tweet-list {
  text-align: left;
}

ul {
  list-style: none;
}
</style>
