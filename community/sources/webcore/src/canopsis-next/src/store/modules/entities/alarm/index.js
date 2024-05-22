import Vue from 'vue';

import { API_ROUTES } from '@/config';

import request, { useRequestCancelling } from '@/services/request';

import i18n from '@/i18n';

import { mapIds } from '@/helpers/array';
import { mergeReactiveChangedProperties } from '@/helpers/vue-base';

import { types as activeWidgetsTypes, getters as activeWidgetGetters } from '../../active-view/active-widgets';

import detailsModule from './details';
import linksModule from './links';

export const types = {
  SET_ALARMS: 'SET_ALARMS',
};

export default {
  namespaced: true,
  modules: {
    details: detailsModule,
    links: linksModule,
  },
  state: {
    alarmsById: {},
    widgets: {},
  },
  getters: {
    getItem: state => id => state.alarmsById[id],

    getList: (state, getters) => ids => ids.map(getters.getItem),

    getListByWidgetId: (state, getters, rootState, rootGetters) => (
      widgetId => getters.getList(rootGetters[activeWidgetGetters.GET_ALL_IDS_BY_WIDGET_ID](widgetId))
    ),

    getMetaByWidgetId: (state, getters, rootState, rootGetters) => (
      widgetId => rootGetters[activeWidgetGetters.GET_META_BY_WIDGET_ID](widgetId)
    ),

    getPendingByWidgetId: (state, getters, rootState, rootGetters) => (
      widgetId => rootGetters[activeWidgetGetters.GET_PENDING_BY_WIDGET_ID](widgetId)
    ),

    getFetchingParamsByWidgetId: (state, getters, rootState, rootGetters) => (
      widgetId => rootGetters[activeWidgetGetters.GET_FETCHING_PARAMS_BY_WIDGET_ID](widgetId)
    ),
  },
  mutations: {
    [types.SET_ALARMS](state, { data }) {
      data.forEach((alarm) => {
        const oldAlarm = state.alarmsById[alarm._id];

        const updatedAlarm = oldAlarm
          ? mergeReactiveChangedProperties(oldAlarm, alarm)
          : alarm;

        Vue.set(state.alarmsById, alarm._id, updatedAlarm);
      });
    },
  },
  actions: {
    fetchComponentAlarmsListWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.componentAlarms, { params });
    },

    fetchResolvedAlarmsListWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.resolvedAlarms, { params });
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
            commit(activeWidgetsTypes.FETCH_LIST, { widgetId, params }, { root: true });
          }

          const { data, meta } = await request.get(API_ROUTES.alarms.list, { params, cancelToken: source.token });

          commit(types.SET_ALARMS, { data });
          commit(activeWidgetsTypes.FETCH_LIST_COMPLETED, {
            widgetId,
            allIds: mapIds(data),
            meta,
          }, { root: true });
        }, `alarms-list-${widgetId}`);
      } catch (err) {
        console.error(err);

        await dispatch('popups/error', { text: i18n.t('errors.default') }, { root: true });

        commit(activeWidgetsTypes.FETCH_LIST_FAILED, { widgetId }, { root: true });
      }
    },

    async fetchItem({ commit }, { id }) {
      const alarm = await request.get(`${API_ROUTES.alarms.list}/${id}`);

      commit(types.SET_ALARMS, { data: [alarm] });

      return alarm;
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

    fetchDisplayNamesWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.alarmDisplayNames, { params });
    },

    fetchExecutionsWithoutStore(context, { id, params } = {}) {
      return request.get(`${API_ROUTES.alarmExecutions}/${id}`, { params });
    },

    updateItemInStore({ commit }, alarm) {
      commit(types.SET_ALARMS, { data: [alarm] });

      return alarm;
    },
  },
};
