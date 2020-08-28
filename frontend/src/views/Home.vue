<template>
  <div class="home">
    <img alt="Vue logo" src="../assets/logo.png">
    <!-- <HelloWorld msg="Welcome to Your Vue.js + TypeScript App"/> -->
    <TweetList />
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import HelloWorld from '@/components/HelloWorld.vue';
import TweetList from '@/components/TweetList.vue';

@Component({
  components: {
    HelloWorld,
    TweetList,
  },
})
export default class Home extends Vue {}
</script>
