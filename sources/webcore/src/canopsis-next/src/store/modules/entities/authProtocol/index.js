import request from '@/services/request';
import { API_ROUTES } from '@/config';

const types = {
  FETCH_CONFIG_STARTED: 'FETCH_CONFIG_STARTED',
  FETCH_CONFIG_COMPLETED: 'FETCH_CONFIG_COMPLETED',
  FETCH_CONFIG_ERROR: 'FETCH_ERROR',
};

export default {
  namespaced: true,
  state: {
    configPending: false,
  },
  mutations: {
    [types.FETCH_CONFIG_STARTED](state) {
      state.configPending = true;
    },
    [types.FETCH_CONFIG_COMPLETED](state) {
      state.configPending = false;
    },
    [types.FETCH_CONFIG_ERROR](state) {
      state.configPending = false;
    },
  },
  actions: {
    async fetchLDAPConfigWithoutStore({ commit }) {
      commit(types.FETCH_CONFIG_STARTED);

      try {
        const { data } = await request.get(API_ROUTES.authProtocols.ldapConfig);
        commit(types.FETCH_CONFIG_COMPLETED);
        return data;
      } catch (err) {
        return commit(types.FETCH_CONFIG_ERROR);
      }
    },

    async updateLDAPConfig(context, { data } = {}) {
      try {
        await request.put(API_ROUTES.authProtocols.ldapConfig, data);
      } catch (err) {
        console.warn(err);
      }
    },

    async fetchCASConfigWithoutStore({ commit }) {
      commit(types.FETCH_CONFIG_STARTED);

      try {
        const { data } = await request.get(API_ROUTES.authProtocols.casConfig);
        commit(types.FETCH_CONFIG_COMPLETED);
        return data;
      } catch (err) {
        return commit(types.FETCH_CONFIG_ERROR);
      }
    },

    async updateCASConfig(context, { data } = {}) {
      try {
        await request.put(API_ROUTES.authProtocols.casConfig, data);
      } catch (err) {
        console.warn(err);
      }
    },
  },
};
