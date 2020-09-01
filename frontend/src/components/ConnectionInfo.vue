<template>
  <section class="connection-info">
    <p>
      websocket connection: 
      <span v-bind:class="connectionOpen">
        {{ connectionOpen }}
      </span>
    </p>
    <button v-on:click="toggleConnection">
      {{ connectionOpen === "closed" ? "open" : "close" }} the websocket
    </button>
    <div v-if="error != null">
      <pre>
        {{ error.toString() }}
      </pre>
    </div>
  </section>
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
/* TODO - extract all of this to a custom component */
button {
	background: linear-gradient(to bottom, #f9f9f9 5%, #e9e9e9 100%);
	background-color: #f9f9f9;
	box-shadow: inset 0px 1px 0px 0px #ffffff;
	border: 1px solid #dcdcdc;
	border-radius: 6px;
	cursor: pointer;
	color: #666666;
	display: inline-block;
	font-family: Arial;
	font-size: 15px;
	font-weight: bold;
	padding: 6px 24px;
	text-decoration: none;
	text-shadow: 0px 1px 0px #ffffff;

  &:hover {
    background: linear-gradient(to bottom, #e9e9e9 5%, #f9f9f9 100%);
    background-color: #e9e9e9;
  }

  &:active {
    position: relative;
    top: 1px;
  }
}

.closed {
  color: red;
}

.open {
  color: green;
}
</style>
