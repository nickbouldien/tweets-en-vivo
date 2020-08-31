<template>
  <div class="connection-info">
    <p>
      websocket connection: 
      <span v-bind:class="connectionOpen">
        {{ connectionOpen }}
      </span>
    </p>
    <button v-on:click="toggleConnection">
      {{ connectionOpen === "closed" ? "open" : "close" }} the websocket connection
    </button>
    <div v-if="error != null">
      <pre>
        {{ error.toString() }}
      </pre>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import Welcome from '@/components/Welcome.vue';
import TweetList from '@/components/TweetList.vue';

@Component({
  name: 'ConnectionInfo',
  components: {
    ConnectionInfo,
  }
})
export default class ConnectionInfo extends Vue {
  @Prop(Boolean) websocketOpen!: boolean;
  @Prop(Error) error!: Error;

  toggleConnection() {
    this.$emit("toggleWebsocket");
  }

  get connectionOpen(): string {
    return Boolean(this.websocketOpen) ? "open" : "closed";
  }
}
</script>

<style scoped lang="scss">
.closed {
  color: red;
}

.open {
  color: green;
}
</style>
