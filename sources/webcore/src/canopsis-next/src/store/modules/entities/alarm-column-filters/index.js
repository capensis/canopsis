import { API_ROUTES } from '@/config';
import request from '@/services/request';

const types = {
  FETCH_ALARM_COLUMN_FILTERS: 'FETCH_ALARM_COLUMN_FILTERS',
  FETCH_ALARM_COLUMN_FILTERS_COMPLETED: 'FETCH_ALARM_COLUMN_FILTERS_COMPLETED',
};

export default {
  namespaced: true,
  state: {
    pending: false,
    alarmColumnFilters: [],
  },
  getters: {
    pending: state => state.pending,
    alarmColumnFilters: state => state.alarmColumnFilters,
  },
  mutations: {
    [types.FETCH_ALARM_COLUMN_FILTERS](state) {
      state.pending = true;
    },
    [types.FETCH_ALARM_COLUMN_FILTERS_COMPLETED](state, { filters }) {
      state.pending = false;
      state.alarmColumnFilters = filters;
    },
  },
  actions: {
    async fetchAlarmColumnFilters({ commit }) {
      try {
        commit(types.FETCH_ALARM_COLUMN_FILTERS);

        const { filters = [] } = await request.get(API_ROUTES.alarmColumnFilters);

        commit(types.FETCH_ALARM_COLUMN_FILTERS_COMPLETED, { filters });
      } catch (err) {
        console.error(err);
      }
    },
  },
};
