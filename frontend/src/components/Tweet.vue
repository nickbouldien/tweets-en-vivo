<template>
  <li class="tweet">
    <a
      :href="tweet.userUrl"
      target="_blank"
      rel="noopener"
    >
      <span class="username">@{{ tweet.authorUsername }}</span>
      <span> · </span>
      <span class="name">{{ tweet.authorName }}</span>
    </a>
    <br />
    <!-- TODO - display the time the tweet was created -->
    <a
      :href="tweet.tweetUrl"
      class="tweet-text"
      target="_blank"
      rel="noopener"
    >
      {{ tweet.text }}
    </a>
  </li>
</template>

<script lang="ts">
import Vue from 'vue';
import { Component, Prop } from 'vue-property-decorator';

import { ITweet } from '@/types';

@Component({
  name: 'Tweet',
  components: {
    Tweet,
  }
})
export default class Tweet extends Vue {
  @Prop({ required: true }) public tweet!: ITweet;
}
</script>

<style lang="scss" scoped>
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
  border-radius: 10px;
  margin: 6px 0;
  padding: 8px;
}

.tweet-text {
  color: #2c3e50;
}

.username {
  color: #42b983;
  font-weight: bold;
}
</style>
