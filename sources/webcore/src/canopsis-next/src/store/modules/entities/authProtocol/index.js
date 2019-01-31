import request from '@/services/request';
import { API_ROUTES } from '@/config';

const types = {
  FETCH_LDAP_CONFIG_STARTED: 'FETCH_LDAP_CONFIG_STARTED',
  FETCH_LDAP_CONFIG_COMPLETED: 'FETCH_LDAP_CONFIG_COMPLETED',
  FETCH_LDAP_CONFIG_ERROR: 'FETCH_LDAP_ERROR',
  UPDATE_LDAP_CONFIG: 'UPDATE_LDAP_CONFIG',
};

export default {
  namespaced: true,
  state: {
    ldapConfigPending: false,
    ldapConfig: {},
  },
  getters: {
    ldapConfig: state => state.ldapConfig,
  },
  mutations: {
    [types.FETCH_LDAP_CONFIG_STARTED](state) {
      state.ldapConfigPending = true;
    },
  },
  actions: {
    async fetchLDAPConfigWithoutStore({ commit }) {
      commit(types.FETCH_LDAP_CONFIG_STARTED);

      const { data } = await request.get(API_ROUTES.authProtocols.ldapConfig);

      return data;
    },
  },
};
