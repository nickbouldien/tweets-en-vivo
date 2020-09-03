<template>
  <div class="welcome">
    <h1>tweets en vivo</h1>
    <p>
      this project will allow you to follow/stream tweets based on criteria or "rules"
    </p>
    <p>
      check out the <router-link to="/about">about page</router-link> to learn more about the project
    </p>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';

@Component
export default class Welcome extends Vue {}
</script>

<style lang="scss" scoped>
a {
  color: #42b983;
}

h1 {
  margin: 12px 0;
}
</style>
