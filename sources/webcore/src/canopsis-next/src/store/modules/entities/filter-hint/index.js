import { API_ROUTES } from '@/config';
import request from '@/services/request';

const types = {
  FETCH_FILTER_HINTS: 'FETCH_FILTER_HINTS',
};

export default {
  namespaced: true,
  state: {
    alarmFilterHints: [],
    entityFilterHints: [],
  },
  getters: {
    alarmFilterHints: state => state.alarmFilterHints,
    entityFilterHints: state => state.entityFilterHints,
  },
  mutations: {
    [types.FETCH_FILTER_HINTS](state, { alarmFilterHints, entityFilterHints }) {
      state.alarmFilterHints = alarmFilterHints;
      state.entityFilterHints = entityFilterHints;
    },
  },
  actions: {
    async fetchFilterHints({ commit }) {
      try {
        const { alarm: alarmFilterHints, entity: entityFilterHints } = await request.get(API_ROUTES.filterHints);

        commit(types.FETCH_FILTER_HINTS, { alarmFilterHints, entityFilterHints });
      } catch (err) {
        console.error(err);
      }
    },
  },
};
