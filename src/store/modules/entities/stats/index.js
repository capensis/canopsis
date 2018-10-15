import i18n from '@/i18n';
import set from 'lodash/set';

import request from '@/services/request';
import { API_ROUTES } from '@/config';

export default {
  namespaced: true,
  actions: {
    async fetchListWithoutStore({ dispatch }, { params, aggregate }) {
      try {
        if (aggregate) {
          Object.keys(params.stats).forEach(stat => set(params.stats[stat], 'aggregate', aggregate));
        }

        return await request.post(API_ROUTES.stats, { ...params });
      } catch (err) {
        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });

        return [];
      }
    },

    async fetchEvolutionWithoutStore({ dispatch }, { params, aggregate }) {
      try {
        if (aggregate) {
          Object.keys(params.stats).forEach(stat => set(params.stats[stat], 'aggregate', aggregate));
        }

        return await request.post(`${API_ROUTES.stats}/evolution`, { ...params });
      } catch (err) {
        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });

        return [];
      }
    },
  },
};
