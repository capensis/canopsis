import i18n from '@/i18n';
import request from '@/services/request';
import { API_ROUTES } from '@/config';

export default {
  namespaced: true,
  actions: {
    /**
     * Action for right creating/editing
     *
     * @param {Object} context
     * @param {Object} data
     * @returns {Promise}
     */
    create(context, { data }) {
      return request.post(API_ROUTES.action, data);
    },

    /**
     * Action for right removing
     *
     * @param {Object} context
     * @param {string} id
     * @returns {Promise}
     */
    remove(context, { id }) {
      return request.delete(`${API_ROUTES.action}/${id}`);
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
