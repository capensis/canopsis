import { merge } from 'lodash';

import request from '@/services/request';

export const DEFAULT_ENTITY_MODULE_TYPES = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

/**
 * Creates a basic CRUD module with create, update, and remove actions.
 *
 * @param {Object} options - Configuration options for the CRUD module.
 * @param {string} options.route - The base API route for the entity.
 * @param {boolean} [options.namespaced = true] - Namespaced flag.
 * @param {boolean} [options.withWithoutStore = false] - Flag to include fetchListWithoutStore action.
 * @param {Object} [module = {}] - Additional module configuration.
 * @returns {Object} An object containing CRUD actions.
 * @property {Function} create - Action to create an entity.
 * @property {Function} update - Action to update an entity by ID.
 * @property {Function} remove - Action to remove an entity by ID.
 * @property {Function} [fetchListWithoutStore] - Action to fetch list without storing (if enabled).
 */
export const createBasicCRUDModule = ({
  route,
  namespaced = true,
  withWithoutStore = false,
} = {}, module = {}) => {
  const moduleActions = {
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

  if (withWithoutStore) {
    moduleActions.fetchListWithoutStore = (context, { params } = {}) => request.get(route, { params });
  }

  return merge({
    namespaced,
    actions: moduleActions,
  }, module);
};

/**
 * Creates a Vuex module for CRUD operations with additional features.
 *
 * @param {Object} options - Configuration options for the CRUD module.
 * @param {Object} [options.types = DEFAULT_ENTITY_MODULE_TYPES] - Custom types for mutations.
 * @param {string} options.route - The base API route for the entity.
 * @param {Function} [options.dataPreparer = d => d?.data] - Function to prepare data from the response.
 * @param {Function} [options.metaPreparer = d => d?.meta] - Function to prepare meta information from the response.
 * @param {boolean} [options.withFetchingParams] - Flag to enable fetching with parameters.
 * @param {boolean} [options.withWithoutStore] - Flag to enable fetching without storing data.
 * @param {Object} [module = {}] - Additional module configuration.
 * @returns {Object} A Vuex module with state, getters, mutations, and actions for CRUD operations.
 */
export const createCRUDModule = ({
  types = DEFAULT_ENTITY_MODULE_TYPES,
  route,
  dataPreparer = d => d?.data,
  metaPreparer = d => d?.meta,
  withFetchingParams,
  withWithoutStore,
}, module = {}) => {
  const moduleState = {
    items: [],
    meta: {},
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
     * Fetches a list of entities from the API.
     *
     * @param {ActionContext} context - The Vuex action context.
     * @param {Object} [params] - Query parameters for the request.
     * @returns {Promise<void>} A promise that resolves when the fetch is complete.
     */
    async fetchList({ commit }, { params } = {}) {
      try {
        commit(types.FETCH_LIST, { params });

        const response = await request.get(route, { params });

        const data = dataPreparer(response);
        const meta = metaPreparer(response);

        commit(types.FETCH_LIST_COMPLETED, { data, meta });
      } catch (err) {
        console.error(err);
        commit(types.FETCH_LIST_FAILED);

        throw err;
      }
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
  }, createBasicCRUDModule({ route }), module);
};
