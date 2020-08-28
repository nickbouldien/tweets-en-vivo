<template>
  <div id="tweet-list">
    <h2>Tweet list</h2> 
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
      <li v-for="tweet in tweets" :key="tweet.id" class="tweet">
        <a :href="'https://twitter.com/' + tweet.authorUsername" target="_blank">
          <span class="username">@{{ tweet.authorUsername }}</span>
          <span class=""> · </span>
          <span class="name">{{ tweet.authorName }}</span>
        </a>
        <br />
        <a :href="'https://twitter.com/random/status/' + tweet.id " class="tweet-text" target="_blank">
          {{ tweet.text }}
        </a>

        <!-- {{ tweet.created_at }} -->
        <!-- {{ tweet.author_id }} -->
      </li>
    </ul>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import { Component, Watch } from 'vue-property-decorator';

import { websocketUrl } from '../config';
import { Tweet, TweetResponse } from '../store/main/state';

@Component
export default class TweetList extends Vue {
  error: Error | null = null;
  tweets: Tweet[] = [];
  socket: WebSocket | null = null;

  mounted() {
    this.error = null;
    this.tweets = [];
    this.socket = null;
  }

  created() {
    this.createWebSocketConnection();
  }

  createWebSocketConnection(): void {
    console.log("Starting connection to WebSocket Server");
    this.socket = new WebSocket(websocketUrl);

    // register websocket event listeners/handlers
    this.socket.onopen = (event: Event) => {
      console.log('socket opened: ', this.socket);
    }

    this.socket.onmessage = (event: MessageEvent) => {
      let tweetResponse: TweetResponse;
      let tweet: Tweet;
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
a {
  text-decoration: none;

  &:hover {
    text-decoration: underline;
  }
}

.name {
  color: #2c3e50;
  font-size: 12px;
}

.tweet {
  background-color: #f6f8fa;
  padding: 8px 4px;
  margin: 4px 0;
}

#tweet-list {
  text-align: left;
}

.tweet-text {
  color: #2c3e50;
}

.username {
  color: #42b983;
  font-weight: bold;
}

ul {
  list-style: none;
}
</style>
