import { set } from 'lodash';

import i18n from '@/i18n';

import request from '@/services/request';
import { API_ROUTES } from '@/config';

export default {
  namespaced: true,
  actions: {
    async fetchItemValuesWithoutStore({ dispatch }, { params }) {
      try {
        const data = await request.post(`${API_ROUTES.stats}`, params);

        return data.values;
      } catch (err) {
        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });

        return [];
      }
    },

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
