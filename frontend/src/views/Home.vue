<template>
  <div class="home">
    <Welcome />
    <hr />
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

<style lang="scss" scoped>
hr {
  border: 0;
  height: 1px;
  background-image: linear-gradient(to right, rgba(0, 0, 0, 0.1), rgba(0, 0, 0, 0.2), rgba(0, 0, 0, 0.1));
  width: 80%;
}
</style>
