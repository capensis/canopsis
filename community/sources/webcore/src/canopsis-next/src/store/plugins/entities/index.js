import Vue from 'vue';
import { get, pick, uniqWith, mergeWith, isEqual } from 'lodash';
import { normalize, denormalize } from 'normalizr';

import { SCHEMA_EMBEDDED_KEY } from '@/config';
import { REQUEST_METHODS } from '@/constants';

import request from '@/services/request';
import schemas from '@/store/schemas';

import { prepareEntitiesToDelete, cloneSchemaWithEmbedded } from './helpers';
import cache from './cache';

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
      return (type, id, withEmbedded = false) => {
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

        const entity = state[type][id];

        if (!entity) {
          return undefined;
        }

        if (!schema.disabledCache && cache.has(entity)) {
          return cache.get(entity);
        }

        const result = denormalize(id, schema, state);

        if (!schema.disabledCache) {
          cache.set(entity, result);
        }

        return result;
      };
    },
    getList(state) {
      return (type, ids = [], withEmbedded = false) => {
        if (typeof type !== 'string') {
          throw new Error('[entities/getList] Missing required argument.');
        }

        if (!state[type] || ids.length === 0) {
          return [];
        }

        let schema = schemas[type];

        if (withEmbedded) {
          schema = cloneSchemaWithEmbedded(schema);
        }

        const { idAttribute, disabledCache } = schema;
        const entities = denormalize(ids, [schema], state)
          .filter(item => !!item);

        if (disabledCache) {
          return entities;
        }

        return entities.map((item) => {
          const entity = state[type][item[idAttribute]];

          if (cache.has(entity)) {
            return cache.get(entity);
          }

          cache.set(entity, item);

          return item;
        });
      };
    },
  },
  mutations: {
    /**
     * @param {Object} state - state of the module
     * @param {Object.<string, Object>} entities - Object of entities
     */
    [internalTypes.ENTITIES_REPLACE](state, entities) {
      cache.clear();

      Object.keys(state).forEach((type) => {
        Vue.set(state, type, entities[type] || {});
      });
    },

    /**
     * @param {Object} state - state of the module
     * @param {Object.<string, Object>} entities - Object of entities
     */
    [internalTypes.ENTITIES_UPDATE](state, entities) {
      Object.keys(entities).forEach((type) => {
        if (!state[type]) {
          Vue.set(state, type, entities[type]);
        } else {
          Object.entries(entities[type]).forEach(([key, entity]) => {
            cache.clearForEntity(state, entity);

            if (state[type][key]) {
              cache.clearForEntity(state, state[type][key]);
            }

            Vue.set(state[type], key, entity);
          });
        }
      });
    },

    /**
     * @param {Object} state - state of the module
     * @param {Object.<string, Object>} entities - Object of entities
     */
    [internalTypes.ENTITIES_MERGE](state, entities) {
      Object.keys(entities).forEach((type) => {
        if (!state[type]) {
          Vue.set(state, type, entities[type]);
        } else {
          Object.entries(entities[type]).forEach(([key, entity]) => {
            const newEntity = mergeWith({}, state[type][key] || {}, entity, (objValue, srcValue) => {
              if (Array.isArray(objValue)) {
                return uniqWith(objValue.concat(srcValue), isEqual);
              }

              return undefined;
            });

            cache.clearForEntity(state, newEntity);

            if (state[type][key]) {
              cache.clearForEntity(state, state[type][key]);
            }

            Vue.set(state[type], key, newEntity);
          });
        }
      });
    },

    /**
     * @param {Object} state - state of the module
     * @param {Object.<string, Object>} entities - Object of entities
     */
    [internalTypes.ENTITIES_DELETE](state, entities) {
      Object.keys(entities).forEach((type) => {
        if (state[type]) {
          Object.entries(entities[type]).forEach(([key, entity]) => {
            cache.delete(entity);

            if (state[type][key]) {
              cache.delete(state[type][key]);
            }

            Vue.delete(state[type], key);
          });
        }
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
     * @property {function} afterCommit - Response data preparer before normalizing
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
        cancelToken,
        method = REQUEST_METHODS.get,
        headers = {},
        params = {},
        dataPreparer = d => d,
        mutationType = internalTypes.ENTITIES_UPDATE,
        afterCommit,
      },
    ) {
      let data;

      const config = { params, headers, cancelToken };

      switch (method.toUpperCase()) {
        case REQUEST_METHODS.get:
          data = await request.get(route, config);
          break;
        case REQUEST_METHODS.post:
          data = await request.post(route, body, config);
          break;
        case REQUEST_METHODS.put:
          data = await request.put(route, body, config);
          break;
        case REQUEST_METHODS.delete:
          data = await request.delete(route, config);
          break;
        default:
          throw new Error(`Invalid method: ${method}`);
      }

      const normalizedData = normalize(dataPreparer(data), schema);
      commit(mutationType, normalizedData.entities);

      if (afterCommit) {
        afterCommit({ data, normalizedData });
      }

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
      return dispatch('sendRequest', { ...config, method: REQUEST_METHODS.post });
    },

    /**
     * Send PUT request
     *
     * @param {VuexActionContext} context
     * @param {EntitiesRequestConfig} config
     * @returns {Promise<EntitiesResponseData>}
     */
    async update({ dispatch }, config) {
      return dispatch('sendRequest', { ...config, method: REQUEST_METHODS.put });
    },

    /**
     * Send DELETE request
     *
     * @param {VuexActionContext} context
     * @param {EntitiesRequestConfig} config
     * @returns {Promise<EntitiesResponseData>}
     */
    async delete({ dispatch }, config) {
      return dispatch('sendRequest', { ...config, method: REQUEST_METHODS.delete });
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
