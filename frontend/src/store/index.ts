import Vue from 'vue';
import Vuex, { StoreOptions } from 'vuex';

import { mainModule } from './main';
import { IState } from './state';

Vue.use(Vuex);

const storeOptions: StoreOptions<IState> = {
  modules: {
    main: mainModule,
  },
};

export const store = new Vuex.Store<IState>(storeOptions);

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
