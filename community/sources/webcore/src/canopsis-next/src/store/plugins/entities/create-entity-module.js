import { merge } from 'lodash';

import request from '@/services/request';

import schemas from '@/store/schemas';

export const DEFAULT_ENTITY_MODULE_TYPES = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default ({
  types = DEFAULT_ENTITY_MODULE_TYPES,
  entityType,
  route,
  dataPreparer,
  withFetchingParams,
  withMeta,
}, module = {}) => {
  const schema = schemas[entityType];

  const moduleState = {
    allIds: [],
    pending: false,
  };

  const moduleGetters = {
    getItemById: (state, getters, rootState, rootGetters) => id => rootGetters['entities/getItem'](entityType, id),

    items: (state, getters, rootState, rootGetters) => rootGetters['entities/getList'](entityType, state.allIds),

    pending: state => state.pending,
  };

  const moduleMutations = {
    [types.FETCH_LIST](state) {
      state.pending = true;
    },
    [types.FETCH_LIST_COMPLETED](state, { allIds }) {
      state.pending = false;
      state.allIds = allIds;
    },
    [types.FETCH_LIST_FAILED](state) {
      state.pending = false;
    },
  };

  const moduleActions = {
    /**
     *
     * @param {ActionContext} context
     * @param {Object} [params] - Query params of request
     * @returns {Promise.<void>}
     */
    async fetchList({ commit, dispatch }, { params } = {}) {
      try {
        commit(types.FETCH_LIST, { params });

        const { normalizedData, data } = await dispatch('entities/fetch', {
          route,
          params,
          dataPreparer,
          schema: [schema],
        }, { root: true });

        commit(types.FETCH_LIST_COMPLETED, {
          ...data,
          allIds: normalizedData.result,
        });
      } catch (err) {
        console.error(err);
        commit(types.FETCH_LIST_FAILED);

        throw err;
      }
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
     * Update entity by id and data
     *
     * @param {ActionContext} context
     * @param {string} id - Id of entity
     * @param {Object} data - Entity data
     * @returns {Promise<AxiosPromise>}
     */
    update(context, { id, data }) {
      return request.put(`${route}/${encodeURIComponent(id)}`, data);
    },

    /**
     * Remove entity by id
     *
     * @param {ActionContext} context
     * @param {string} id - Id of entity
     * @returns {Promise<AxiosPromise>}
     */
    remove(context, { id } = {}) {
      return request.delete(`${route}/${encodeURIComponent(id)}`);
    },
  };

  if (withFetchingParams) {
    moduleMutations[types.FETCH_LIST] = (state, { params } = {}) => {
      state.pending = true;
      state.fetchingParams = params;
    };

    moduleActions.fetchListWithPreviousParams = ({ dispatch, state }) => dispatch('fetchList', {
      params: state.fetchingParams,
    });
  }

  if (withMeta) {
    moduleState.meta = {};

    moduleMutations[types.FETCH_LIST_COMPLETED] = (state, { allIds, meta }) => {
      state.allIds = allIds;
      state.meta = meta;
      state.pending = false;
    };

    moduleGetters.meta = state => state.meta;
  }

  return merge({
    namespaced: true,
    state: moduleState,
    getters: moduleGetters,
    mutations: moduleMutations,
    actions: moduleActions,
  }, module);
};
