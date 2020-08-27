import { getStoreAccessors } from 'typesafe-vuex';

import { MainState, Tweet } from './state';
import { State } from '../state';

export const mutations = {
  setError(state: MainState, payload: Error) {
    state.error = payload;
  },
  setTweet(state: MainState, tweet: Tweet) {
    state.tweets.push(tweet);
  },
  SOCKET_CONNECT(state: MainState) {
    console.log("SOCKET_CONNECT");
    // state.connected = true;
  },
  SOCKET_DISCONNECT(state: MainState) {
    console.log("SOCKET_DISCONNECT");
    // state.connected = false;
  },
  SOCKET_MESSAGE(state: MainState, message: any) {
    console.log("SOCKET_MESSAGE ", message);
    // state.message = message
  },
  // SOCKET_HELLO_WORLD(state: MainState, message) {
  //   state.message = message
  // },
  SOCKET_ERROR(state: MainState, message) {
    console.log("SOCKET_ERROR ", message);
    // state.error = message.error
  },

};

const { commit } = getStoreAccessors<MainState | any, State>('');

export const commitSetError = commit(mutations.setError);
export const commitSetTweet = commit(mutations.setTweet);
