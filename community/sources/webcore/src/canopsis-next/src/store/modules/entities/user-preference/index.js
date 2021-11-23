import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import request from '@/services/request';

import { userPreferenceSchema } from '@/store/schemas';

export default {
  namespaced: true,
  state: {
    pending: false,
  },
  getters: {
    getItemByWidget: (state, getters, rootState, rootGetters) => (widget) => {
      const userPreference = rootGetters['entities/getItem'](ENTITIES_TYPES.userPreference, widget._id);

      if (!userPreference) {
        return {
          widget: widget._id,
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
     * This action creates user preference
     *
     * @param {function} dispatch
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
  },
};
