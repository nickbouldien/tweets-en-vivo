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
      <li v-for="tweet in tweets" :key="tweet.id">
        <!-- {{ tweet.text }} -->
        tweet received
      </li>
    </ul>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import { Component } from 'vue-property-decorator';

import { Tweet } from '../store/main/state';

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
    console.log("called created ")
    this.createWebSocketConnection();
  }

  createWebSocketConnection(): void {
    console.log("Starting connection to WebSocket Server");
    // TODO - get from env variables
    const port = 5000;
    const path = "/ws";
    const websocketUrl = `ws://localhost:${port}${path}`;
    this.socket = new WebSocket(websocketUrl);

    // register websocket event listeners/handlers
    this.socket.onopen = (event: Event) => {
      console.log('socket opened: ', this.socket);
    }

    this.socket.onmessage = (event: MessageEvent) => {
      let tweet: Tweet;
      try {
        console.log("event.data: ", event.data);
        tweet = JSON.parse(event.data);
        console.log("tweet: ", tweet);
        // this.updateMessages(msg.data);
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
      console.warn("socket closed");
      this.socket = null;
    }
  }

  get connectionOpen(): string {
    console.log("connectionOpen ", this.socket);
    return this.socket ? "open" : "closed";
  }
}
</script>

<style scoped>
#tweet-list {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
