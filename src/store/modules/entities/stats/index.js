import { API_ROUTES } from '@/config';
import { statSchema } from '@/store/schemas';
import { ENTITIES_TYPES } from '@/constants';

export const types = {
  FETCH_STATS: 'FETCH_STATS',
  FETCH_STATS_COMPLETED: 'FETCH_STATS_COMPLETED',
  FETCH_STATS_FAILED: 'FETCH_STATS_FAILED',
};

export default {
  namespaced: true,
  state: {
    pending: false,
    error: {},
  },
  getters: {
    getItem: (state, getters, rootState, rootGetters) => id =>
      rootGetters['entities/getItem'](ENTITIES_TYPES.stat, id),
    getList: (state, getters, rootState, rootGetters) => rootGetters['entities/getList'](ENTITIES_TYPES.stat, state.allIds),
  },
  mutations: {
    [types.FETCH_STATS](state) {
      state.pending = true;
    },
    [types.FETCH_STATS_COMPLETED](state, { allIds }) {
      state.allIds = allIds;
    },
    [types.FETCH_STATS_FAILED](state, error) {
      state.error = error;
    },
  },
  actions: {
    async fetchStats({ dispatch, commit }) {
      // Example params
      const exampleParams = {
        mfilter: {
          type: 'component',
        },
        tstop: 1534716000,
        duration: '2d',
        stats: {
          'Taux de disponibilité': {
            stat: 'state_rate',
            parameters: {
              states: [0, 1, 2],
            },
            trend: true,
            sla: '>= 0.99',
          },
          'Taux de performance': {
            stat: 'state_rate',
            parameters: {
              states: [0],
            },
            trend: true,
            sla: '>= 0.95',
          },
          'Durée d\'indisponibilité': {
            stat: 'time_in_state',
            parameters: {
              states: [3],
            },
            trend: true,
          },
          'Nombre d\'indisponibilités': {
            stat: 'alarms_created',
            parameters: {
              states: [3],
            },
            trend: true,
            sla: '<= 100',
          },
        },
      };

      commit(types.FETCH_STATS);

      try {
        const { normalizedData } = await dispatch('entities/fetch', {
          route: API_ROUTES.stats,
          schema: [statSchema],
          body: exampleParams,
          method: 'POST',
          dataPreparer: d => d.values,
        }, { root: true });

        commit(types.FETCH_STATS_COMPLETED, {
          allIds: normalizedData.result,
        });
      } catch (err) {
        commit(types.FETCH_STATS_FAILED, err);
        console.error(err);
      }
    },
  },
};
