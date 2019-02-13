import { merge } from 'lodash';

import request from '@/services/request';

import schemas from '@/store/schemas';

export default ({
  types,
  entityType,
  route,
  dataPreparer,
}, module = {}) => {
  const schema = schemas[entityType];

  return merge({
    namespaced: true,
    state: {
      allIds: [],
      pending: false,
    },
    getters: {
      items: (state, getters, rootState, rootGetters) =>
        rootGetters['entities/getList'](entityType, state.allIds),
      pending: state => state.pending,
    },
    mutations: {
      [types.FETCH_LIST](state, { params }) {
        state.pending = true;
        state.fetchingParams = params;
      },
      [types.FETCH_LIST_COMPLETED](state, { allIds }) {
        state.pending = false;
        state.allIds = allIds;
      },
      [types.FETCH_LIST_FAILED](state) {
        state.pending = false;
      },
    },
    actions: {
      async fetchList({ commit, dispatch }, { params } = {}) {
        try {
          commit(types.FETCH_LIST, { params });

          const { normalizedData } = await dispatch('entities/fetch', {
            route,
            params,
            dataPreparer,
            schema: [schema],
          }, { root: true });

          commit(types.FETCH_LIST_COMPLETED, {
            allIds: normalizedData.result,
          });
        } catch (err) {
          console.error(err);
          commit(types.FETCH_LIST_FAILED);
        }
      },

      fetchListWithPreviousParams({ dispatch, state }) {
        return dispatch('fetchList', {
          params: state.fetchingParams,
        });
      },

      create(context, { data }) {
        return request.post(route, data);
      },

      edit(context, { id, data }) {
        return request.put(`${route}/${id}`, data);
      },

      async remove({ dispatch }, { id } = {}) {
        try {
          await request.delete(`${route}/${id}`);
          await dispatch('entities/removeFromStore', {
            id,
            type: entityType,
          }, { root: true });
        } catch (err) {
          console.error(err);
        }
      },
    },
  }, module);
};
