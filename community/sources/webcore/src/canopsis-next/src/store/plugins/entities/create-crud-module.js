import { merge } from 'lodash';

import request from '@/services/request';

export const DEFAULT_ENTITY_MODULE_TYPES = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export const createCRUDModule = ({
  types = DEFAULT_ENTITY_MODULE_TYPES,
  route,
  withFetchingParams,
  withWithoutStore,
}, module = {}) => {
  const moduleState = {
    items: [],
    meta: [],
    pending: false,
  };

  const moduleGetters = {
    items: state => state.items,
    meta: state => state.meta,
    pending: state => state.pending,
  };

  const moduleMutations = {
    [types.FETCH_LIST](state) {
      state.pending = true;
    },
    [types.FETCH_LIST_COMPLETED](state, { data, meta }) {
      state.items = data;
      state.meta = meta;
      state.pending = false;
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
    async fetchList({ commit }, { params } = {}) {
      try {
        commit(types.FETCH_LIST, { params });

        const { data, meta } = await request.get(route, { params });

        commit(types.FETCH_LIST_COMPLETED, { data, meta });
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

  if (withWithoutStore) {
    moduleActions.fetchListWithoutStore = (context, options) => request.get(route, options);
  }

  return merge({
    namespaced: true,
    state: moduleState,
    getters: moduleGetters,
    mutations: moduleMutations,
    actions: moduleActions,
  }, module);
};
