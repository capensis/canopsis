import request from '@/services/request';
import { API_ROUTES } from '@/config';

export const types = {
  FETCH_VERSION_COMPLETED: 'FETCH_VERSION_COMPLETED',
};

export default {
  namespaced: true,
  state: {
    version: '',
  },
  getters: {
    version: state => state.version,
  },
  mutations: {
    [types.FETCH_VERSION_COMPLETED](state, version) {
      state.version = version;
    },
  },
  actions: {
    async fetchVersion({ commit }) {
      const { version } = await request.get(API_ROUTES.version);
      commit(types.FETCH_VERSION_COMPLETED, version);
    },
  },
};
