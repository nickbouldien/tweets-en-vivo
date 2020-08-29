// import { actions } from './actions';
import { getters } from './getters';
import { mutations } from './mutations';
import { IMainState } from './state';

const defaultState: IMainState = {
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
