import Vue from 'vue';
import { get } from 'lodash';

const COUNTER = {
  total: 2937,
  total_active: 2582,
  snooze: 0,
  ack: 36,
  ticket: 14,
  pbehavior_active: 355,
};

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default {
  namespaced: true,
  state: {
    widgets: {},
  },
  getters: {
    getListByWidgetId: state => widgetId => get(state.widgets[widgetId], 'counters', []),
    getPendingByWidgetId: state => widgetId => get(state.widgets[widgetId], 'pending', []),
  },
  mutations: {
    [types.FETCH_LIST](state, { widgetId }) {
      Vue.setSeveral(state.widgets, widgetId, { pending: true });
    },
    [types.FETCH_LIST_COMPLETED](state, { widgetId, counters }) {
      Vue.setSeveral(state.widgets, widgetId, { counters, pending: false });
    },
    [types.FETCH_LIST_FAILED](state, { widgetId }) {
      Vue.setSeveral(state.widgets, widgetId, { pending: false });
    },
  },
  actions: {
    async fetchList({ commit }, { widgetId, filters = [] }) {
      try {
        commit(types.FETCH_LIST, { widgetId });

        const promises = filters.map(async () => ({ data: [COUNTER] }));
        const responses = await Promise.all(promises);

        commit(types.FETCH_LIST_COMPLETED, {
          widgetId,

          counters: responses.map(({ data: [counter] }) => counter),
        });
      } catch (err) {
        commit(types.FETCH_LIST_FAILED, { widgetId });
      }
    },
  },
};
