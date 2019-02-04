import omit from 'lodash/omit';

import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';
import { userPreferenceSchema } from '@/store/schemas';
import { generateUserPreferenceByWidgetAndUser } from '@/helpers/entities';
import request from '@/services/request';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
  SET_ACTIVE_FILTER: 'SET_ACTIVE_FILTER',
};

export default {
  namespaced: true,
  getters: {
    getItemByWidget: (state, getters, rootState, rootGetters) => (widget) => {
      const currentUser = rootGetters['auth/currentUser'];
      const id = `${widget._id}_${currentUser.crecord_name}`;
      const userPreference = rootGetters['entities/getItem'](ENTITIES_TYPES.userPreference, id);

      if (!userPreference) {
        return generateUserPreferenceByWidgetAndUser(widget, currentUser);
      }

      return userPreference;
    },
  },
  mutations: {
    [types.FETCH_LIST]: state => state.pending = true,
    [types.FETCH_LIST_COMPLETED]: (state) => {
      state.pending = false;
    },
    [types.FETCH_LIST_FAILED]: (state) => {
      state.pending = false;
    },
  },
  actions: {
    /**
     * This action fetches user preferences list
     *
     * @param {function} commit
     * @param {function} dispatch
     * @param {Object} params
     */
    async fetchList({ commit, dispatch }, { params }) {
      try {
        commit(types.FETCH_LIST);

        await dispatch('entities/fetch', {
          route: API_ROUTES.userPreferences,
          schema: [userPreferenceSchema],
          params,
          dataPreparer: d => d.data,
        }, { root: true });

        commit(types.FETCH_LIST_COMPLETED);
      } catch (err) {
        commit(types.FETCH_LIST_FAILED);
        console.warn(err);
      }
    },

    /**
     * This action fetches user preference item by widget id
     *
     * @param {function} dispatch
     * @param {function} rootGetters
     * @param {string|number} widgetId
     */
    fetchItemByWidgetId({ dispatch, rootGetters }, { widgetId }) {
      const currentUser = rootGetters['auth/currentUser'];

      return dispatch('fetchList', {
        params: {
          limit: 1,
          filter: {
            crecord_name: currentUser.crecord_name,
            widget_id: widgetId,
            _id: `${widgetId}_${currentUser.crecord_name}`,
          },
        },
      });
    },

    /**
     * This action fetches user preference item by widget id without store
     *
     * @param {function} rootGetters
     * @param {string|number} widgetId
     */
    async fetchItemByWidgetIdWithoutStore({ rootGetters }, { widgetId }) {
      const currentUser = rootGetters['auth/currentUser'];
      const params = {
        limit: 1,
        filter: {
          crecord_name: currentUser.crecord_name,
          widget_id: widgetId,
          _id: `${widgetId}_${currentUser.crecord_name}`,
        },
      };

      const { data: [item] } = await request.get(API_ROUTES.userPreferences, { params });

      return item;
    },

    /**
     * This action creates user preference
     *
     * @param {function} dispatch
     * @param {Object} userPreference
     */
    async create({ dispatch }, { userPreference }) {
      try {
        const body = omit(userPreference, ['crecord_creation_time', 'crecord_write_time', 'enable']);

        await dispatch('entities/update', {
          route: API_ROUTES.userPreferences,
          schema: userPreferenceSchema,
          body: JSON.stringify(body),
          dataPreparer: d => d.data[0],
        }, { root: true });
      } catch (err) {
        console.warn(err);
      }
    },
  },
};
