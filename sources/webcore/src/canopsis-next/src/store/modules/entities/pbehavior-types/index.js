import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';
import { pbehaviorTypesSchema } from '@/store/schemas';

import { createEntityModule } from '@/store/plugins/entities';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default createEntityModule({
  types,
  route: API_ROUTES.planning.types,
  entityType: ENTITIES_TYPES.pbehaviorTypes,
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
          route: API_ROUTES.planning.types,
          schema: [pbehaviorTypesSchema],
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
  },
});
