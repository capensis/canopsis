import { API_ROUTES } from '@/config';
import request from '@/services/request';

const types = {
  FETCH_FILTER_HINTS: 'FETCH_FILTER_HINTS',
  FETCH_FILTER_HINTS_COMPLETED: 'FETCH_FILTER_HINTS_COMPLETED',
};

export default {
  namespaced: true,
  state: {
    pending: false,
    alarmFilterHints: [],
    entityFilterHints: [],
  },
  getters: {
    pending: state => state.pending,
    alarmFilterHints: state => state.alarmFilterHints,
    entityFilterHints: state => state.entityFilterHints,
  },
  mutations: {
    [types.FETCH_FILTER_HINTS](state) {
      state.pending = true;
    },
    [types.FETCH_FILTER_HINTS_COMPLETED](state, { alarmFilterHints, entityFilterHints }) {
      state.pending = false;
      state.alarmFilterHints = alarmFilterHints;
      state.entityFilterHints = entityFilterHints;
    },
  },
  actions: {
    async fetchFilterHints({ commit }) {
      try {
        commit(types.FETCH_FILTER_HINTS);

        const { alarm: alarmFilterHints, entity: entityFilterHints } = await request.get(API_ROUTES.filterHints);

        commit(types.FETCH_FILTER_HINTS_COMPLETED, { alarmFilterHints, entityFilterHints });
      } catch (err) {
        console.error(err);
      }
    },
  },
};
