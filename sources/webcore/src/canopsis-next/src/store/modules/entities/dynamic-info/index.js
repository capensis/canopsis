import i18n from '@/i18n';
import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';
import { dynamicInfoSchema } from '@/store/schemas';

import { createEntityModule } from '@/store/plugins/entities';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default createEntityModule({
  types,
  route: API_ROUTES.dynamicInfo,
  entityType: ENTITIES_TYPES.dynamicInfo,
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
          route: API_ROUTES.dynamicInfo,
          params,
          schema: [dynamicInfoSchema],
          dataPreparer: d => d.rules,
        }, { root: true });

        commit(types.FETCH_LIST_COMPLETED, {
          allIds: normalizedData.result,
          meta: {
            total: data.count,
          },
        });
      } catch (err) {
        commit(types.FETCH_LIST_FAILED);

        await dispatch('popups/error', { text: i18n.t('errors.default') }, { root: true });
      }
    },
  },
});
