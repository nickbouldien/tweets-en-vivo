<template>
  <div class="home">
    <Welcome />
    <TweetList />
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import Welcome from '@/components/Welcome.vue';
import TweetList from '@/components/TweetList.vue';

@Component({
  components: {
    TweetList,
    Welcome,
  },
})
export default class Home extends Vue {}
</script>
