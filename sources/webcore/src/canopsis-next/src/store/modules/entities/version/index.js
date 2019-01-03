import request from '@/services/request';
import i18n from '@/i18n';
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
    async fetchVersion({ commit, dispatch }) {
      try {
        const { version } = await request.get(API_ROUTES.version);
        commit(types.FETCH_VERSION_COMPLETED, version);
      } catch (err) {
        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.versionNotFound') }, { root: true });
        console.warn(err);
      }
    },
  },
};
