import i18n from '@/i18n';
import request from '@/services/request';
import { API_ROUTES } from '@/config';

export default {
  namespaced: true,
  actions: {
    /**
     * Action for right creating/editing
     *
     * @param {Function} dispatch
     * @param {Object} data
     * @returns void
     */
    async create({ dispatch }, { data }) {
      try {
        await request.post(API_ROUTES.action, data);
        await dispatch('popup/add', { type: 'success', text: i18n.t('success.default') }, { root: true });
      } catch (err) {
        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });
      }
    },

    /**
     * Fetch rights list by params without store
     *
     * @param {Function} commit
     * @param {Function} dispatch
     * @param {Object} params
     * @returns {Object}
     */
    async fetchListWithoutStore({ dispatch }, { params }) {
      try {
        return request.get(API_ROUTES.action, { params });
      } catch (err) {
        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });

        return { data: [], total: 0 };
      }
    },
  },
};
