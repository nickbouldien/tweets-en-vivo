import Vue from 'vue'
import VueSocketIO from 'vue-socket.io'

import App from './App.vue'
import router from './router'
import store from './store'

Vue.config.productionTip = false

// TODO - get from env variables
// const port = 5000;
// const path = "/ws";
// const websocketUrl = `ws://localhost:${port}${path}`;

// Vue.use(new VueSocketIO({
//   debug: true,
//   connection: websocketUrl,
//   vuex: {
//       store,
//       actionPrefix: 'SOCKET_',
//       mutationPrefix: 'SOCKET_'
//   },
// }));


new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app');
