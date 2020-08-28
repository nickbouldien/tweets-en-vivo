import { getStoreAccessors } from 'typesafe-vuex';

import { MainState } from './state';
import { State } from '../state';

export const getters = {
  error: (state: MainState) => state.error,
  tweets: (state: MainState) => state.tweets,
};

const { read } = getStoreAccessors<MainState, State>('');

export const readError = read(getters.error);
export const readTweets = read(getters.tweets);
