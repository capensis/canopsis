import Vue from 'vue';
import get from 'lodash/get';
import omit from 'lodash/omit';
import uniq from 'lodash/uniq';
import pickBy from 'lodash/pickBy';
import mergeWith from 'lodash/mergeWith';
import isEmpty from 'lodash/isEmpty';
import { normalize, denormalize } from 'normalizr';

import request from '@/services/request';
import schemas from '@/store/schemas';
import { prepareEntitiesToDelete } from '@/helpers/store';

const entitiesModuleName = 'entities';

const internalTypes = {
  ENTITIES_UPDATE: 'ENTITIES_UPDATE',
  ENTITIES_MERGE: 'ENTITIES_MERGE',
  ENTITIES_DELETE: 'ENTITIES_DELETE',
};

const usingEntitiesCount = {};

const entitiesModule = {
  namespaced: true,
  state: {
    entities: {},
    entitiesUsingCount: {},
  },
  getters: {
    getItem(state) {
      return (type, id) => {
        if (typeof type !== 'string') {
          throw new Error('[entities/getItem] Missing required argument.');
        }

        if (!state.entities[type] || !id) {
          return null;
        }

        return denormalize(id, schemas[type], state.entities);
      };
    },
    getList(state) {
      return (type, ids = []) => {
        if (typeof type !== 'string') {
          throw new Error('[entities/getList] Missing required argument.');
        }

        if (!state.entities[type] || ids.length === 0) {
          return [];
        }

        const result = denormalize(ids, [schemas[type]], state.entities);

        return result.filter(v => !!v);
      };
    },
  },
  mutations: {
    /**
     * @param {Object} state - state of the module
     * @param {Object.<string, Object>} entities - Object of entities
     */
    [internalTypes.ENTITIES_UPDATE](state, entities) {
      // TODO: Put usingCount update and filter state.entities by it
      Object.keys(entities).forEach((type) => {
        Vue.set(state.entities, type, {
          ...(state.entities[type] || {}),
          ...entities[type],
        });
      });
    },

    /**
     * @param {Object} state - state of the module
     * @param {Object.<string, Object>} entities - Object of entities
     */
    [internalTypes.ENTITIES_MERGE](state, entities) {
      // TODO: Put usingCount update and filter state.entities by it
      Object.keys(entities).forEach((type) => {
        Vue.set(
          state.entities, type, mergeWith({}, state.entities[type] || {}, entities[type]),
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
        Vue.set(state.entities, key, omit(state.entities[key], Object.keys(entities[key])));
      });
    },
  },
  actions: {
    sweep({ dispatch }, { type } = {}) {
      let entitiesForDeletion = {};

      if (type) {
        entitiesForDeletion[type] = pickBy(usingEntitiesCount[type], value => value <= 0);
      } else {
        entitiesForDeletion = Object.keys(usingEntitiesCount).reduce((acc, localType) => {
          const items = pickBy(usingEntitiesCount[localType], value => value <= 0);

          if (!isEmpty(items)) {
            acc[localType] = items;
          }

          return acc;
        }, {});
      }

      Object.keys(entitiesForDeletion).map(localType => Object.keys(entitiesForDeletion[localType])
        .map(id => dispatch('removeFromStore', { id, type: localType })));
    },

    mark(context, { type, usingCount = {} } = {}) {
      if (!usingEntitiesCount[type]) {
        usingEntitiesCount[type] = {};
      }

      Object.keys(usingCount).forEach((key) => {
        if (!usingEntitiesCount[type][key]) {
          usingEntitiesCount[type][key] = 0;
        }

        usingEntitiesCount[type][key] += usingCount[key];
      });
    },

    async sendRequest(
      { commit },
      {
        route,
        schema,
        body,
        headers = {},
        dataPreparer = d => d,
        params = {},
        method = 'GET',
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
    async fetch({ dispatch }, config) {
      const newConfig = omit(config, ['isPost']);

      if (config.isPost) {
        newConfig.method = 'POST';
      }

      return dispatch('sendRequest', newConfig);
    },
    async create({ dispatch }, config) {
      return dispatch('sendRequest', { ...config, method: 'POST' });
    },
    async update({ dispatch }, config) {
      return dispatch('sendRequest', { ...config, method: 'PUT' });
    },
    async delete({ dispatch }, config) {
      return dispatch('sendRequest', { ...config, method: 'DELETE' });
    },

    /**
     * Remove entity by id and type from store
     */
    removeFromStore({ commit, getters, state }, { id, type }) {
      const data = getters.getItem(type, id);
      const parents = get(data, '_embedded.parents', []);

      const {
        entitiesToMerge,
        entitiesToDelete,
      } = prepareEntitiesToDelete({ type, data });

      parents.forEach((parent) => {
        const parentEntity = state.entities[parent.type][parent.id];

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
  ENTITIES_DELETE: `${entitiesModuleName}/${internalTypes.ENTITIES_DELETE}`,
};

export default (store) => {
  store.registerModule(entitiesModuleName, entitiesModule);
};
