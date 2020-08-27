import Vue from 'vue';
import Vuex, { StoreOptions } from 'vuex';

import { mainModule } from './main';
import { State } from './state';

Vue.use(Vuex);

// function createWebSocketPlugin (socket) {
//   return store => {
//     socket.on('data', data => {
//       store.commit('receiveData', data)
//     })
//     store.subscribe(mutation => {
//       if (mutation.type === 'UPDATE_DATA') {
//         socket.emit('update', mutation.payload)
//       }
//     })
//   }
// }

// const plugin = createWebSocketPlugin(socket)

const storeOptions: StoreOptions<State> = {
  modules: {
    main: mainModule,
  },
};

export const store = new Vuex.Store<State>(storeOptions);

console.log("store: ", store);

// const store = new Vuex.Store({
//   state: {
//   },
//   mutations: {
//   },
//   actions: {
//   },
//   modules: {
//   },
//   // plugins: [plugin],
// });

export default store;
