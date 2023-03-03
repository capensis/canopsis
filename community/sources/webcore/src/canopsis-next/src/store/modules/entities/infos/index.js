import { MAX_LIMIT } from '@/constants';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default {
  namespaced: true,
  state: {
    alarmInfos: [],
    alarmInfosRules: [],
    entityInfos: [],
    pending: false,
  },
  getters: {
    alarmInfos: state => state.alarmInfos ?? [],
    alarmInfosRules: state => state.alarmInfosRules ?? [],
    entityInfos: state => state.entityInfos ?? [],
    pending: state => state.pending ?? false,
  },
  mutations: {
    [types.FETCH_LIST]: (state) => {
      state.pending = true;
    },

    [types.FETCH_LIST_COMPLETED]: (state, { alarmInfos = [], alarmInfosRules = [], entityInfos = [] }) => {
      state.alarmInfos = alarmInfos;
      state.alarmInfosRules = alarmInfosRules;
      state.entityInfos = entityInfos;
      state.pending = false;
    },

    [types.FETCH_LIST_FAILED]: (state) => {
      state.pending = false;
    },
  },
  actions: {
    async fetch({ dispatch, commit }, { params = { limit: MAX_LIMIT } } = {}) {
      try {
        commit(types.FETCH_LIST);

        const requests = [
          dispatch('dynamicInfo/fetchInfosKeysWithoutStore', { params }, { root: true }),
          dispatch('service/fetchInfosKeysWithoutStore', { params }, { root: true }),
          dispatch('dynamicInfo/fetchListWithoutStore', { params }, { root: true }),
        ];

        const [
          { data: alarmInfos } = {},
          { data: entityInfos } = {},
          { data: alarmInfosRules } = {},
        ] = await Promise.all(requests);

        commit(types.FETCH_LIST_COMPLETED, { alarmInfos, entityInfos, alarmInfosRules });
      } catch (err) {
        console.error(err);

        commit(types.FETCH_LIST_FAILED);
      }
    },
  },
};
