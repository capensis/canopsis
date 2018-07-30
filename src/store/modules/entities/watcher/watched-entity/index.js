import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';
import i18n from '@/i18n/index';
import { watchedEntitiesSchema } from '@/store/schemas/index';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default {
  namespaced: true,
  state: {
    allIds: [],
    meta: {},
    pending: false,
    fetchingParams: {},
  },
  getters: {
    allIds: state => state.allIds,
    items: (state, getters, rootState, rootGetters) => rootGetters['entities/getList'](ENTITIES_TYPES.watchedEntity, state.allIds),
    pending: state => state.pending,
    meta: state => state.meta,
    item: (state, getters, rootState, rootGetters) => id => rootGetters['entities/getItem'](ENTITIES_TYPES.watchedEntity, id),
  },
  mutations: {
    [types.FETCH_LIST](state, { params }) {
      state.pending = true;
      state.fetchingParams = params;
    },
    [types.FETCH_LIST_COMPLETED](state, { allIds, meta }) {
      state.allIds = allIds;
      state.meta = meta;
      state.pending = false;
    },
    [types.FETCH_LIST_FAILED](state) {
      state.pending = false;
    },
  },
  actions: {
    async fetchList({ dispatch, commit }, { watcherId, params }) {
      try {
        commit(types.FETCH_LIST, { params });

        const { normalizedData } = await dispatch('entities/fetch', {
          route: `${API_ROUTES.watchers}/${watcherId}`,
          schema: [watchedEntitiesSchema],
          params,
          dataPreparer: d => d,
        }, { root: true });

        commit(types.FETCH_LIST_COMPLETED, { allIds: normalizedData.result, meta: {} });
      } catch (e) {
        commit(types.FETCH_LIST_FAILED);
        console.error(e);
        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });
      }
    },
  },
};

