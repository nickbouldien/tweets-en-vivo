// import { actions } from './actions';
import { getters } from './getters';
import { mutations } from './mutations';
import { MainState } from './state';

const defaultState: MainState = {
  error: null,
  tweets: [],
  websocket: null,
};

export const mainModule = {
  state: defaultState,
  mutations,
  // actions,
  getters,
};
