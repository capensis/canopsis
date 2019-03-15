import Vue from 'vue';
import { get, omit, pick, uniq, mergeWith } from 'lodash';
import { normalize, denormalize } from 'normalizr';

import request from '@/services/request';
import schemas from '@/store/schemas';
import { prepareEntitiesToDelete, cloneSchemaWithEmbedded } from '@/helpers/store';
import { SCHEMA_EMBEDDED_KEY } from '@/config';

const entitiesModuleName = 'entities';

const internalTypes = {
  ENTITIES_UPDATE: 'ENTITIES_UPDATE',
  ENTITIES_MERGE: 'ENTITIES_MERGE',
  ENTITIES_REPLACE: 'ENTITIES_REPLACE',
  ENTITIES_DELETE: 'ENTITIES_DELETE',
};

let registeredGetters = [];

export const entitiesModule = {
  namespaced: true,
  getters: {
    getItem(state) {
      return (type, id, withEmbedded) => {
        let schema = schemas[type];

        if (typeof type !== 'string') {
          throw new Error('[entities/getItem] Missing required argument.');
        }

        if (withEmbedded) {
          schema = cloneSchemaWithEmbedded(schema);
        }

        if (!state[type] || !id) {
          return null;
        }

        return denormalize(id, schema, state);
      };
    },
    getList(state) {
      return (type, ids = []) => {
        if (typeof type !== 'string') {
          throw new Error('[entities/getList] Missing required argument.');
        }

        if (!state[type] || ids.length === 0) {
          return [];
        }

        const result = denormalize(ids, [schemas[type]], state);

        return result.filter(v => !!v);
      };
    },
  },
  mutations: {
    /**
     * @param {Object} state - state of the module
     * @param {Object.<string, Object>} entities - Object of entities
     */
    [internalTypes.ENTITIES_REPLACE](state, entities) {
      Object.keys(entities).forEach((type) => {
        Vue.set(state, type, entities[type]);
      });
    },

    /**
     * @param {Object} state - state of the module
     * @param {Object.<string, Object>} entities - Object of entities
     */
    [internalTypes.ENTITIES_UPDATE](state, entities) {
      Object.keys(entities).forEach((type) => {
        Vue.set(state, type, {
          ...(state[type] || {}),
          ...entities[type],
        });
      });
    },

    /**
     * @param {Object} state - state of the module
     * @param {Object.<string, Object>} entities - Object of entities
     */
    [internalTypes.ENTITIES_MERGE](state, entities) {
      Object.keys(entities).forEach((type) => {
        Vue.set(
          state, type, mergeWith({}, state[type] || {}, entities[type]),
          (objValue, srcValue) => {
            if (Array.isArray(objValue)) {
              return uniq(objValue.concat(srcValue));
            }

            return undefined;
          },
        );
      });
    },

    /**
     * @param {Object} state - state of the module
     * @param {Object.<string, Object>} entities - Object of entities
     */
    [internalTypes.ENTITIES_DELETE](state, entities) {
      Object.keys(entities).forEach((key) => {
        Vue.set(state, key, omit(state[key], Object.keys(entities[key])));
      });
    },
  },
  actions: {
    /**
     * Register getterObject
     *
     * @param {VuexActionContext} context
     * @param {Object} getterObject - getter object for registration
     * @param {function} getterObject.getDependencies - Method for getting component dependencies
     * @param {Vue.component} getterObject.instance - Instance of component
     */
    registerGetter(context, getterObject) {
      registeredGetters.push(getterObject);
    },

    /**
     * Unregister getterObject by instance
     *
     * @param {VuexActionContext} context
     * @param {Vue.component} instance - Instance of component
     */
    unregisterGetter(context, instance) {
      registeredGetters = registeredGetters.filter(getterObject => getterObject.instance !== instance);
    },

    /**
     * Sweep unregistered entities
     *
     * @param {VuexActionContext} context
     */
    sweep({ commit, state }) {
      const entities = registeredGetters
        .reduce((acc, { getDependencies }) => acc.concat(getDependencies()), [])
        .reduce((acc, { type, ids }) => {
          acc[type] = { ...(acc[type]), ...pick(state[type], ids) };

          return acc;
        }, {});

      commit(internalTypes.ENTITIES_REPLACE, entities);
    },

    /**
     * @typedef {Object} EntitiesRequestConfig
     * @property {string} route - Route of resource
     * @property {Object} schema - Schema for the resource
     * @property {Object|string} body - Request body
     * @property {string} method - Method of the request
     * @property {Object} headers - Request headers
     * @property {Object} params - Request query params
     * @property {function} dataPreparer - Response data preparer before normalizing
     * @property {string} mutationType - Mutation type after normalization
     */

    /**
     * @typedef {Object} EntitiesResponseData
     * @property {any} data - Response data
     * @property {Object} normalizedData - Response normalized data
     * @property {Object} normalizedData.entities - Object based entities
     * @property {Array} normalizedData.result - Ids of main entities
     */

    /**
     * Send request by our request service
     *
     * @param {VuexActionContext} context
     * @param {EntitiesRequestConfig} config
     * @returns {Promise<EntitiesResponseData>}
     */
    async sendRequest(
      { commit },
      {
        route,
        schema,
        body,
        method = 'GET',
        headers = {},
        params = {},
        dataPreparer = d => d,
        mutationType = internalTypes.ENTITIES_UPDATE,
      },
    ) {
      let data;

      switch (method) {
        case 'GET':
          data = await request.get(route, { params, headers });
          break;
        case 'POST':
          data = await request.post(route, body, { params, headers });
          break;
        case 'PUT':
          data = await request.put(route, body, { params, headers });
          break;
        case 'DELETE':
          data = await request.delete(route, { params, headers });
          break;
        default:
          throw new Error(`Invalid method: ${method}`);
      }

      const normalizedData = normalize(dataPreparer(data), schema);
      commit(mutationType, normalizedData.entities);

      return { data, normalizedData };
    },

    /**
     * Send GET request
     *
     * @param {VuexActionContext} context
     * @param {EntitiesRequestConfig} config
     * @returns {Promise<EntitiesResponseData>}
     */
    async fetch({ dispatch }, config) {
      return dispatch('sendRequest', config);
    },

    /**
     * Send POST request
     *
     * @param {VuexActionContext} context
     * @param {EntitiesRequestConfig} config
     * @returns {Promise<EntitiesResponseData>}
     */
    async create({ dispatch }, config) {
      return dispatch('sendRequest', { ...config, method: 'POST' });
    },

    /**
     * Send PUT request
     *
     * @param {VuexActionContext} context
     * @param {EntitiesRequestConfig} config
     * @returns {Promise<EntitiesResponseData>}
     */
    async update({ dispatch }, config) {
      return dispatch('sendRequest', { ...config, method: 'PUT' });
    },

    /**
     * Send DELETE request
     *
     * @param {VuexActionContext} context
     * @param {EntitiesRequestConfig} config
     * @returns {Promise<EntitiesResponseData>}
     */
    async delete({ dispatch }, config) {
      return dispatch('sendRequest', { ...config, method: 'DELETE' });
    },

    /**
     * Remove entity by id and type from store
     *
     * @param {VuexActionContext} context
     * @param {Object} payload
     * @param {string|number} payload.id - Id of entity for deletion
     * @param {string} payload.type - Type of entity for deletion
     */
    removeFromStore({ commit, getters, state }, { id, type }) {
      const data = getters.getItem(type, id, true);
      const parents = get(data, `${SCHEMA_EMBEDDED_KEY}.parents`, []);

      const {
        entitiesToMerge,
        entitiesToDelete,
      } = prepareEntitiesToDelete({ type, data });

      parents.forEach((parent) => {
        const parentEntity = state[parent.type][parent.id];

        if (!entitiesToMerge[parent.type]) {
          entitiesToMerge[parent.type] = {};
        }

        entitiesToMerge[parent.type][parent.id] = {
          ...parentEntity,
          [parent.key]: parentEntity[parent.key].filter(v => v !== id),
        };
      });

      commit(internalTypes.ENTITIES_UPDATE, entitiesToMerge);
      commit(internalTypes.ENTITIES_DELETE, entitiesToDelete);
    },
  },
};

export const types = {
  ENTITIES_UPDATE: `${entitiesModuleName}/${internalTypes.ENTITIES_UPDATE}`,
  ENTITIES_MERGE: `${entitiesModuleName}/${internalTypes.ENTITIES_MERGE}`,
  ENTITIES_REPLACE: `${entitiesModuleName}/${internalTypes.ENTITIES_REPLACE}`,
  ENTITIES_DELETE: `${entitiesModuleName}/${internalTypes.ENTITIES_DELETE}`,
};

export { default as createEntityModule } from './create-entity-module';

export default (store) => {
  store.registerModule(entitiesModuleName, entitiesModule);
};
