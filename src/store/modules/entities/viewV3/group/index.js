import request from '@/services/request';
import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';
import i18n from '@/i18n';
import { groupSchema } from '@/store/schemas';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default {
  namespaced: true,
  state: {
    allIds: [],
    pending: false,
    fetchingParams: {},
  },
  getters: {
    allIds: state => state.allIds,
    items: (state, getters, rootState, rootGetters) => rootGetters['entities/getList'](ENTITIES_TYPES.group, state.allIds),
    pending: state => state.pending,
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
    async create(context, params = {}) {
      try {
        const route = API_ROUTES.viewV3.groups;
        return await request.post(route, params);
      } catch (err) {
        console.warn(err);
        return undefined;
      }
    },
    async fetchList({ commit, dispatch }, { name } = {}) {
      try {
        let route = API_ROUTES.viewV3.groups;
        if (name) {
          route += name;
        }
        const { normalizedData } = await dispatch('entities/fetch', {
          route,
          schema: [groupSchema],
          dataPreparer: d => d,
        }, { root: true });

        commit(types.FETCH_LIST_COMPLETED, {
          allIds: normalizedData.result,
        });
      } catch (err) {
        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });
        commit(types.FETCH_LIST_FAILED);
      }
    },
  },
};
