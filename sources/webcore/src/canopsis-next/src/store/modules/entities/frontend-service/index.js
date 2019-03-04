import request from '@/services/request';
import i18n from '@/i18n';
import { API_ROUTES } from '@/config';

export const types = {
  FETCH_ITEM_COMPLETED: 'FETCH_ITEM_COMPLETED',
};

export default {
  namespaced: true,
  state: {
    item: {},
  },
  getters: {
    item: state => state.item,
  },
  mutations: {
    [types.FETCH_ITEM_COMPLETED](state, item) {
      state.item = item;
    },
  },
  actions: {
    async fetch({ commit, dispatch }) {
      try {
        const { data: [item] } = await request.get(API_ROUTES.frontendService);

        commit(types.FETCH_ITEM_COMPLETED, item);
      } catch (err) {
        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });
        console.warn(err);
      }
    },

    async update({ commit, dispatch }, { data }) {
      try {
        const { data: [item] } = await request.get(API_ROUTES.frontendService, data);

        commit(types.FETCH_ITEM_COMPLETED, item);
      } catch (err) {
        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });
        console.warn(err);
      }
    },
  },
};
