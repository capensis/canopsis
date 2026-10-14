import { normalize } from 'normalizr';

import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import request from '@/services/request';

import { userPreferenceSchema } from '@/store/schemas';
import { types as entitiesTypes } from '@/store/plugins/entities';

export default {
  namespaced: true,
  state: {
    pending: false,
  },
  getters: {
    getItemByWidgetId: (state, getters, rootState, rootGetters) => (widgetId) => {
      const userPreference = rootGetters['entities/getItem'](ENTITIES_TYPES.userPreference, widgetId);

      if (!userPreference) {
        return {
          widget: widgetId,
          content: {},
        };
      }

      return userPreference;
    },
  },
  actions: {
    /**
     * This action fetches user preferences list
     *
     * @param {function} dispatch
     * @param {string} id
     */
    async fetchItem({ dispatch }, { id }) {
      try {
        await dispatch('entities/fetch', {
          route: `${API_ROUTES.userPreferences}/${id}`,
          schema: userPreferenceSchema,
        }, { root: true });
      } catch (err) {
        console.warn(err);
      }
    },

    /**
     * This action fetches user preference item by widget id without store
     *
     * @param {VuexActionContext} context
     * @param {string} id
     */
    fetchItemWithoutStore(context, { id }) {
      return request.get(`${API_ROUTES.userPreferences}/${id}`);
    },

    /**
     * This action updates user preference
     *
     * @param {Function} dispatch
     * @param {Object} data
     */
    async update({ dispatch }, { data }) {
      try {
        await dispatch('entities/update', {
          route: API_ROUTES.userPreferences,
          schema: userPreferenceSchema,
          body: data,
        }, { root: true });
      } catch (err) {
        console.warn(err);
      }
    },

    /**
     * This action updates user preference but only in the store (without request)
     *
     * @param {Function} commit
     * @param {Object} data
     */
    updateLocal({ commit }, { data }) {
      const { entities } = normalize(data, userPreferenceSchema);

      commit(entitiesTypes.ENTITIES_UPDATE, entities, { root: true });
    },
  },
};
