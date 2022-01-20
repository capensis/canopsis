import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import request from '@/services/request';

import { createEntityModule } from '@/store/plugins/entities';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',

  FETCH_ACTIVE_LIST: 'FETCH_ACTIVE_LIST',
  FETCH_ACTIVE_LIST_COMPLETED: 'FETCH_ACTIVE_LIST_COMPLETED',
  FETCH_ACTIVE_LIST_FAILED: 'FETCH_ACTIVE_LIST_FAILED',
};

export default createEntityModule({
  types,
  route: API_ROUTES.broadcastMessage.list,
  entityType: ENTITIES_TYPES.broadcastMessage,
  dataPreparer: d => d.data,
  withFetchingParams: true,
  withMeta: true,
}, {
  state: {
    activeMessagesIds: [],
    activeMessagesPending: false,
  },
  getters: {
    activeMessages: (state, getters, rootState, rootGetters) => rootGetters['entities/getList'](
      ENTITIES_TYPES.broadcastMessage,
      state.activeMessagesIds,
    ),

    activeMessagesPending: state => state.activeMessagesPending,
  },
  mutations: {
    [types.FETCH_ACTIVE_LIST](state) {
      state.activeMessagesPending = true;
    },
    [types.FETCH_ACTIVE_LIST_COMPLETED](state, { allIds }) {
      state.activeMessagesIds = allIds;
    },
    [types.FETCH_ACTIVE_LIST_FAILED](state) {
      state.activeMessagesPending = false;
    },
  },
  actions: {
    async fetchActiveListWithoutStore() {
      return request.get(API_ROUTES.broadcastMessage.activeList);
    },
  },
});
