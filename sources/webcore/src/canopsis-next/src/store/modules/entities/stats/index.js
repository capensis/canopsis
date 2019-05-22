import i18n from '@/i18n';

import request from '@/services/request';
import { API_ROUTES } from '@/config';
import { STATS_LISTS_NAMES } from '@/constants';

export default {
  namespaced: true,
  actions: {
    async fetchListWithoutStore({ dispatch }, { params = {} } = {}) {
      try {
        return request.post(API_ROUTES.stats, params);
      } catch (err) {
        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });

        return {
          aggregations: {},
          values: [],
        };
      }
    },

    async fetchEvolutionWithoutStore({ dispatch }, { params = {} } = {}) {
      try {
        return request.post(`${API_ROUTES.stats}/evolution`, params);
      } catch (err) {
        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });

        return {
          aggregations: {},
          values: [],
        };
      }
    },

    async fetchSpecialStatsListsWithoutStore(
      { dispatch },
      {
        listName = STATS_LISTS_NAMES.stateIntervals,
        params = {},
      } = {},
    ) {
      try {
        return request.post(`${API_ROUTES.stats}/lists/${listName}`, params);
      } catch (err) {
        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });

        return [];
      }
    },
  },
};
