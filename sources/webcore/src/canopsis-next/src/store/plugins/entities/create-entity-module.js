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
      /**
       *
       * @param {ActionContext} context
       * @param {Object} [params={}] - Query params of request
       * @returns {Promise.<void>}
       */
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

      /**
       * Fetch list with previous fetching params
       *
       * @param {ActionContext} context
       * @returns {Promise<AxiosPromise>}
       */
      fetchListWithPreviousParams({ dispatch, state }) {
        return dispatch('fetchList', {
          params: state.fetchingParams,
        });
      },

      /**
       * Create entity by data
       *
       * @param {ActionContext} context
       * @param {Object} data - Entity data
       * @returns {Promise<AxiosPromise>}
       */
      create(context, { data }) {
        return request.post(route, data);
      },

      /**
       * Edit entity by id and data
       *
       * @param {ActionContext} context
       * @param {string} id - Id of entity
       * @param {Object} data - Entity data
       * @returns {Promise<AxiosPromise>}
       */
      edit(context, { id, data }) {
        return request.put(`${route}/${id}`, data);
      },

      /**
       * Remove entity by id
       *
       * @param {ActionContext} context
       * @param {string} id - Id of entity
       * @returns {Promise<AxiosPromise>}
       */
      async remove(context, { id } = {}) {
        return request.delete(`${route}/${id}`);
      },
    },
  }, module);
};
