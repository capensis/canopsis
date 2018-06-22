// LIBS
import qs from 'qs';
import i18n from '@/i18n';
// SERVICES
import request from '@/services/request';
// OTHERS
import { API_ROUTES } from '@/config';

export default {
  namespaced: true,
  actions: {
    async create({ dispatch, context }, { data }) {
      try {
        await request.post(API_ROUTES.event, qs.stringify({ event: JSON.stringify(data) }), {
          headers: { 'content-type': 'application/x-www-form-urlencoded' },
        });
        await dispatch('popup/add', { type: 'success', text: i18n.t('success.default') }, { root: true });
      } catch (e) {
        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });
        console.warn(e);
      }
    },
  },
};
