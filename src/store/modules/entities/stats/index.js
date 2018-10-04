import i18n from '@/i18n';

import request from '@/services/request';
import { API_ROUTES } from '@/config';

export default {
  namespaced: true,
  actions: {
    async fetchListWithoutStore({ dispatch }, { params }) {
      try {
        const data = await request.post(API_ROUTES.stats, { ...params });

        return data.values;
      } catch (err) {
        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });

        return [];
      }
    },
  },
};
