import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import request from '@/services/request';

import { createEntityModule } from '@/store/plugins/entities';
import { pbehaviorDatesExceptionsSchema } from '@/store/schemas';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default createEntityModule({
  types,
  route: API_ROUTES.planning.datesExceptions,
  entityType: ENTITIES_TYPES.pbehaviorDatesExceptions,
  withFetchingParams: true,
}, {
  state: {
    meta: {},
  },
  getters: {
    meta: state => state.meta,
  },
  mutations: {
    [types.FETCH_LIST_COMPLETED](state, { allIds, meta }) {
      state.allIds = allIds;
      state.meta = meta;
      state.pending = false;
    },
  },
  actions: {
    async fetchList({ commit, dispatch }, { params } = {}) {
      try {
        commit(types.FETCH_LIST, { params });

        const { normalizedData, data } = await dispatch('entities/fetch', {
          params,
          route: API_ROUTES.planning.datesExceptions,
          schema: [pbehaviorDatesExceptionsSchema],
          dataPreparer: d => d.data,
        }, { root: true });

        commit(types.FETCH_LIST_COMPLETED, {
          allIds: normalizedData.result,
          meta: data.meta,
        });
      } catch (err) {
        console.error(err);
        commit(types.FETCH_LIST_FAILED);

        throw err;
      }
    },

    fetchListWithoutStore(context, { params }) {
      return request.get(API_ROUTES.planning.datesExceptions, { params });
    },
  },
});
