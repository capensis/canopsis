import Vue from 'vue';
import { get } from 'lodash';

import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import request, { useRequestCancelling } from '@/services/request';

import i18n from '@/i18n';

import { alarmSchema } from '@/store/schemas';
import { ENTITIES_MUTATION_TYPES } from '@/store/plugins/entities';

import detailsModule from './details';
import linksModule from './links';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',

  DOWNLOAD_LIST_COMPLETED: 'DOWNLOAD_LIST_COMPLETED',
};

export default {
  namespaced: true,
  modules: {
    details: detailsModule,
    links: linksModule,
  },
  state: {
    widgets: {},
  },
  getters: {
    getListByWidgetId: (state, getters, rootState, rootGetters) => widgetId => rootGetters['entities/getList'](
      ENTITIES_TYPES.alarm,
      get(state.widgets[widgetId], 'allIds', []),
    ),

    getMetaByWidgetId: state => widgetId => get(state.widgets[widgetId], 'meta', {}),
    getPendingByWidgetId: state => widgetId => get(state.widgets[widgetId], 'pending', false),
    getFetchingParamsByWidgetId: state => widgetId => get(state.widgets[widgetId], 'fetchingParams'),

    getItem: (state, getters, rootState, rootGetters) => id => rootGetters['entities/getItem'](
      ENTITIES_TYPES.alarm,
      id,
    ),
    getList: (state, getters, rootState, rootGetters) => ids => rootGetters['entities/getList'](
      ENTITIES_TYPES.alarm,
      ids,
    ),
  },
  mutations: {
    [types.FETCH_LIST](state, { widgetId, params }) {
      Vue.setSeveral(state.widgets, widgetId, { pending: true, fetchingParams: params });
    },
    [types.FETCH_LIST_COMPLETED](state, { widgetId, allIds, meta }) {
      Vue.setSeveral(state.widgets, widgetId, { allIds, meta, pending: false });
    },
    [types.FETCH_LIST_FAILED](state, { widgetId }) {
      Vue.setSeveral(state.widgets, widgetId, { pending: false });
    },
  },
  actions: {
    fetchComponentAlarmsListWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.componentAlarms, { params });
    },

    fetchResolvedAlarmsListWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.resolvedAlarms, { params });
    },

    fetchManualMetaAlarmsListWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.manualMetaAlarm, { params });
    },

    fetchOpenAlarmsListWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.openAlarms, { params });
    },

    fetchListWithoutStore(context, { params }) {
      return request.get(API_ROUTES.alarms.list, { params });
    },

    async fetchList({ commit, dispatch }, { widgetId, params, withoutPending } = {}) {
      try {
        await useRequestCancelling(async (source) => {
          if (!withoutPending) {
            commit(types.FETCH_LIST, { widgetId, params });
          }

          await dispatch('entities/fetch', {
            route: API_ROUTES.alarms.list,
            schema: [alarmSchema],
            params,
            cancelToken: source.token,
            dataPreparer: d => d.data,
            mutationType: ENTITIES_MUTATION_TYPES.smartUpdate,
            afterCommit: ({ normalizedData, data }) => {
              commit(types.FETCH_LIST_COMPLETED, {
                widgetId,
                allIds: normalizedData.result,
                meta: data.meta,
              });
            },
          }, { root: true });
        }, `alarms-list-${widgetId}`);
      } catch (err) {
        console.error(err);

        await dispatch('popups/error', { text: i18n.t('errors.default') }, { root: true });

        commit(types.FETCH_LIST_FAILED, { widgetId });
      }
    },

    fetchItem({ dispatch }, { id }) {
      return dispatch('entities/fetch', {
        route: `${API_ROUTES.alarms.list}/${id}`,
        schema: alarmSchema,
      }, { root: true });
    },

    fetchItemWithoutStore(context, { id }) {
      return request.get(`${API_ROUTES.alarms.list}/${id}`);
    },

    fetchAlarmsListExport(context, { params, id }) {
      return request.get(`${API_ROUTES.alarmListExport}/${id}`, { params });
    },

    async createAlarmsListExport(context, { data = {} }) {
      return request.post(API_ROUTES.alarmListExport, data);
    },

    createAlarmAckEvent(context, { id, data }) {
      return request.put(`${API_ROUTES.alarms.list}/${id}/ack`, data);
    },

    bulkCreateAlarmAckEvent(context, { data }) {
      return request.put(`${API_ROUTES.alarms.bulkList}/ack`, data);
    },

    createAlarmAckremoveEvent(context, { id, data }) {
      return request.put(`${API_ROUTES.alarms.list}/${id}/ackremove`, data);
    },

    bulkCreateAlarmAckremoveEvent(context, { data }) {
      return request.put(`${API_ROUTES.alarms.bulkList}/ackremove`, data);
    },

    createAlarmSnoozeEvent(context, { id, data }) {
      return request.put(`${API_ROUTES.alarms.list}/${id}/snooze`, data);
    },

    bulkCreateAlarmSnoozeEvent(context, { data }) {
      return request.put(`${API_ROUTES.alarms.bulkList}/snooze`, data);
    },

    createAlarmAssocticketEvent(context, { id, data }) {
      return request.put(`${API_ROUTES.alarms.list}/${id}/assocticket`, data);
    },

    bulkCreateAlarmAssocticketEvent(context, { data }) {
      return request.put(`${API_ROUTES.alarms.bulkList}/assocticket`, data);
    },

    createAlarmCommentEvent(context, { id, data }) {
      return request.put(`${API_ROUTES.alarms.list}/${id}/comment`, data);
    },

    bulkCreateAlarmCommentEvent(context, { data }) {
      return request.put(`${API_ROUTES.alarms.bulkList}/comment`, data);
    },

    createAlarmCancelEvent(context, { id, data }) {
      return request.put(`${API_ROUTES.alarms.list}/${id}/cancel`, data);
    },

    bulkCreateAlarmCancelEvent(context, { data }) {
      return request.put(`${API_ROUTES.alarms.bulkList}/cancel`, data);
    },

    createAlarmUncancelEvent(context, { id, data }) {
      return request.put(`${API_ROUTES.alarms.list}/${id}/uncancel`, data);
    },

    bulkCreateAlarmUnCancelEvent(context, { data }) {
      return request.put(`${API_ROUTES.alarms.bulkList}/uncancel`, data);
    },

    createAlarmChangestateEvent(context, { id, data }) {
      return request.put(`${API_ROUTES.alarms.list}/${id}/changestate`, data);
    },

    bulkCreateAlarmChangestateEvent(context, { data }) {
      return request.put(`${API_ROUTES.alarms.bulkList}/changestate`, data);
    },

    addBookmarkToAlarm(context, { id }) {
      return request.put(`${API_ROUTES.alarms.list}/${id}/bookmark`);
    },

    removeBookmarkFromAlarm(context, { id }) {
      return request.delete(`${API_ROUTES.alarms.list}/${id}/bookmark`);
    },

    updateItemInStore({ dispatch }, alarm) {
      return dispatch('entities/addToStore', {
        schema: alarmSchema,
        data: alarm,
      }, { root: true });
    },
  },
};
