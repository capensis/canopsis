import Vue from 'vue';
import {
  get,
  pick,
  uniqWith,
  mergeWith,
  isEqual,
} from 'lodash';
import { normalize, denormalize } from 'normalizr';

import { SCHEMA_EMBEDDED_KEY } from '@/config';
import { REQUEST_METHODS } from '@/constants';

import request from '@/services/request';

import schemas from '@/store/schemas';

import { mergeChangedProperties } from '@/helpers/collection';

import { prepareEntitiesToDelete, cloneSchemaWithEmbedded } from './helpers';
import cache from './cache';

const entitiesModuleName = 'entities';

export const ENTITIES_MUTATION_TYPES = {
  update: 'ENTITIES_UPDATE',
  smartUpdate: 'ENTITIES_SMART_UPDATE',
  merge: 'ENTITIES_MERGE',
  replace: 'ENTITIES_REPLACE',
  delete: 'ENTITIES_DELETE',
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
    [ENTITIES_MUTATION_TYPES.replace](state, entities) {
      cache.clear();

      Object.keys(state).forEach((type) => {
        Vue.set(state, type, entities[type] || {});
      });
    },

    /**
     * @param {Object} state - state of the module
     * @param {Object.<string, Object>} entities - Object of entities
     */
    [ENTITIES_MUTATION_TYPES.update](state, entities) {
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
    [ENTITIES_MUTATION_TYPES.smartUpdate](state, entities) {
      Object.keys(entities).forEach((type) => {
        if (!state[type]) {
          Vue.set(state, type, entities[type]);
        } else {
          Object.entries(entities[type]).forEach(([key, entity]) => {
            const oldEntity = state[type][key];

            const data = oldEntity
              ? mergeChangedProperties(oldEntity, entity)
              : entity;

            Vue.set(state[type], key, data);
          });
        }
      });
    },

    /**
     * @param {Object} state - state of the module
     * @param {Object.<string, Object>} entities - Object of entities
     */
    [ENTITIES_MUTATION_TYPES.merge](state, entities) {
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
    [ENTITIES_MUTATION_TYPES.delete](state, entities) {
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

      commit(ENTITIES_MUTATION_TYPES.replace, entities);
    },

    /**
     * @typedef {Object} EntitiesNormalizationConfig
     * @property {Object} schema - Schema for the resource
     * @property {Object} data - Data for normalization
     * @property {string} mutationType - Mutation type after normalization
     */

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
      { dispatch },
      {
        route,
        schema,
        body,
        cancelToken,
        method = REQUEST_METHODS.get,
        headers = {},
        params = {},
        dataPreparer = d => d,
        mutationType = ENTITIES_MUTATION_TYPES.update,
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

      const normalizedData = await dispatch('addToStore', {
        schema,
        data: dataPreparer(data),
        mutationType,
      });

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
    fetch({ dispatch }, config) {
      return dispatch('sendRequest', config);
    },

    /**
     * Send POST request
     *
     * @param {VuexActionContext} context
     * @param {EntitiesRequestConfig} config
     * @returns {Promise<EntitiesResponseData>}
     */
    create({ dispatch }, config) {
      return dispatch('sendRequest', { ...config, method: REQUEST_METHODS.post });
    },

    /**
     * Send PUT request
     *
     * @param {VuexActionContext} context
     * @param {EntitiesRequestConfig} config
     * @returns {Promise<EntitiesResponseData>}
     */
    update({ dispatch }, config) {
      return dispatch('sendRequest', { ...config, method: REQUEST_METHODS.put });
    },

    /**
     * Send DELETE request
     *
     * @param {VuexActionContext} context
     * @param {EntitiesRequestConfig} config
     * @returns {Promise<EntitiesResponseData>}
     */
    delete({ dispatch }, config) {
      return dispatch('sendRequest', { ...config, method: REQUEST_METHODS.delete });
    },

    /**
     *
     * @param {VuexActionContext} context
     * @param schema
     * @param data
     * @param mutationType
     * @returns {*}
     */
    addToStore({ commit }, { schema, data, mutationType = ENTITIES_MUTATION_TYPES.update }) {
      const normalizedData = normalize(data, schema);

      commit(mutationType, normalizedData.entities);

      return normalizedData;
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

      commit(ENTITIES_MUTATION_TYPES.update, entitiesToMerge);
      commit(ENTITIES_MUTATION_TYPES.delete, entitiesToDelete);
    },
  },
};

export const types = {
  ENTITIES_UPDATE: `${entitiesModuleName}/${ENTITIES_MUTATION_TYPES.update}`,
  ENTITIES_MERGE: `${entitiesModuleName}/${ENTITIES_MUTATION_TYPES.merge}`,
  ENTITIES_REPLACE: `${entitiesModuleName}/${ENTITIES_MUTATION_TYPES.replace}`,
  ENTITIES_DELETE: `${entitiesModuleName}/${ENTITIES_MUTATION_TYPES.delete}`,
};

export { default as createEntityModule } from './create-entity-module';
export { createWidgetModule } from './create-widget-module';
export { createCRUDModule } from './create-crud-module';

export default (store) => {
  store.registerModule(entitiesModuleName, entitiesModule);
};
