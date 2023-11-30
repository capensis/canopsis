import Vue from 'vue';
import { get } from 'lodash';

import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import request from '@/services/request';

import { alarmDetailsSchema } from '@/store/schemas';

import { convertAlarmDetailsQueryToRequest } from '@/helpers/entities/alarm/query';
import { generateAlarmDetailsId, getAlarmDetailsDataPreparer } from '@/helpers/entities/alarm/list';

export const types = {
  FETCH_ITEM: 'FETCH_ITEM',
  FETCH_ITEM_COMPLETED: 'FETCH_ITEM_COMPLETED',

  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',

  UPDATE_QUERY: 'UPDATE_QUERY',
  REMOVE_QUERY: 'REMOVE_QUERY',
};

export default {
  namespaced: true,
  state: {
    widgets: {},
  },
  getters: {
    getItem: (state, getters, rootState, rootGetters) => (widgetId, id) => rootGetters['entities/getItem'](
      ENTITIES_TYPES.alarmDetails,
      generateAlarmDetailsId(id, widgetId),
    ),

    getPending: state => (widgetId, id) => (
      get(state.widgets, [widgetId, id, 'pending'], false)
    ),

    getQuery: state => (widgetId, id) => (
      get(state.widgets, [widgetId, id, 'query'], {})
    ),

    getQueries: state => widgetId => (
      Object.values(state.widgets[widgetId] ?? {}).map(({ query }) => query)
    ),
  },
  mutations: {
    [types.FETCH_ITEM]: (state, { widgetId, id }) => {
      if (!state.widgets[widgetId]) {
        Vue.set(state.widgets, widgetId, { [id]: { pending: true, query: {} } });
      } else if (state.widgets[widgetId][id]) {
        Vue.set(state.widgets[widgetId][id], 'pending', true);
      } else {
        Vue.set(state.widgets[widgetId], id, { pending: true, query: {} });
      }
    },

    [types.FETCH_ITEM_COMPLETED]: (state, { widgetId, id }) => {
      Vue.set(state.widgets[widgetId][id], 'pending', false);
    },

    [types.FETCH_LIST]: (state, { widgetId }) => {
      if (state.widgets[widgetId]) {
        Object.values(state.widgets[widgetId]).forEach(item => Vue.set(item, 'pending', true));
      }
    },

    [types.FETCH_LIST_COMPLETED]: (state, { widgetId }) => {
      if (state.widgets[widgetId]) {
        Object.values(state.widgets[widgetId]).forEach(item => Vue.set(item, 'pending', false));
      }
    },

    [types.UPDATE_QUERY]: (state, { widgetId, id, query }) => {
      if (!state.widgets[widgetId]) {
        Vue.set(state.widgets, widgetId, { [id]: { pending: false, query } });
      } else if (state.widgets[widgetId][id]) {
        Vue.set(state.widgets[widgetId][id], 'query', query);
      } else {
        Vue.set(state.widgets[widgetId], id, { pending: false, query });
      }
    },

    [types.REMOVE_QUERY]: (state, { widgetId, id }) => {
      Vue.delete(state.widgets[widgetId], id);
    },
  },
  actions: {
    /**
     * Fetch alarm details for widget by special query
     *
     * @param {Function} dispatch
     * @param {Function} commit
     * @param {string} widgetId
     * @param {string} id
     * @param {Object} query
     * @returns {Promise<void>}
     */
    async fetchItem({ dispatch, commit }, { widgetId, id, query }) {
      try {
        commit(types.FETCH_ITEM, { widgetId, id });

        await dispatch('entities/create', {
          route: API_ROUTES.alarmDetails,
          schema: [alarmDetailsSchema],
          body: [query],
          dataPreparer: getAlarmDetailsDataPreparer(widgetId),
        }, { root: true });
      } catch (err) {
        console.error(err);
      } finally {
        commit(types.FETCH_ITEM_COMPLETED, { widgetId, id });
      }
    },

    /**
     * Fetch alarms details list for widget (only for all opened expand panel in the widget)
     *
     * @param {Function} dispatch
     * @param {Function} commit
     * @param {Object} state
     * @param {string} widgetId
     * @returns {Promise<void>}
     */
    async fetchList({ dispatch, commit, getters }, { widgetId }) {
      try {
        const queries = getters.getQueries(widgetId);

        if (!queries.length) {
          return;
        }

        commit(types.FETCH_LIST, { widgetId });

        await dispatch('entities/create', {
          route: API_ROUTES.alarmDetails,
          schema: [alarmDetailsSchema],
          body: queries.map(convertAlarmDetailsQueryToRequest),
          dataPreparer: getAlarmDetailsDataPreparer(widgetId),
        }, { root: true });
      } catch (err) {
        console.error(err);
      } finally {
        commit(types.FETCH_LIST_COMPLETED, { widgetId });
      }
    },

    /**
     * Fetch alarms details list without store
     *
     * @param {Object} context
     * @param {Object} params
     * @returns {Promise<void>}
     */
    async fetchListWithoutStore(context, { params } = {}) {
      return request.post(API_ROUTES.alarmDetails, params);
    },

    /**
     * Update query for special alarm details
     *
     * @param {Function} commit
     * @param {string} widgetId
     * @param {string} id
     * @param {Object} query
     */
    updateQuery({ commit }, { widgetId, id, query }) {
      commit(types.UPDATE_QUERY, { widgetId, id, query });
    },

    /**
     * Remove whole query object for special alarm details
     *
     * @param {Function} commit
     * @param {string} widgetId
     * @param {string} id
     */
    removeQuery({ commit }, { widgetId, id }) {
      commit(types.REMOVE_QUERY, { widgetId, id });
    },

    /**
     * Update alarm details item in store
     *
     * @param {Function} dispatch
     * @param {Object} alarmDetails
     * @returns {*}
     */
    updateItemInStore({ dispatch }, alarmDetails) {
      return dispatch('entities/addToStore', {
        schema: alarmDetailsSchema,
        data: alarmDetails,
      }, { root: true });
    },
  },
};
