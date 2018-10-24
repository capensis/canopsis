import i18n from '@/i18n';
import qs from 'qs';

import request from '@/services/request';
import { API_ROUTES } from '@/config';

export default {
  namespaced: true,
  actions: {
    async edit({ dispatch }, { user } = {}) {
      try {
        await request.post(API_ROUTES.user, qs.stringify({ user: JSON.stringify(user) }), {
          headers: { 'content-type': 'application/x-www-form-urlencoded' },
        });
      } catch (err) {
        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });
      }
    },
  },
};
