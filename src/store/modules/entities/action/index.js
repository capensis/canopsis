import i18n from '@/i18n';
import request from '@/services/request';
import { API_ROUTES } from '@/config';

export default {
  namespaced: true,
  actions: {
    async fetchListWithoutStore({ dispatch }, { params }) {
      try {
        return await request.get(API_ROUTES.action, { params });
      } catch (err) {
        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });

        return { data: [], total: 0 };
      }
    },
  },
};
