import { getStoreAccessors } from 'typesafe-vuex';

import { IMainState } from './state';
import { IState } from '../state';

export const getters = {
  error: (state: IMainState) => state.error,
  tweets: (state: IMainState) => state.tweets,
};

const { read } = getStoreAccessors<IMainState, IState>('');

export const readError = read(getters.error);
export const readTweets = read(getters.tweets);
