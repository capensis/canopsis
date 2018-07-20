import Vue from 'vue';
import omit from 'lodash/omit';
import uniq from 'lodash/uniq';
import mergeWith from 'lodash/mergeWith';
import get from 'lodash/get';
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

const entitiesModule = {
  namespaced: true,
  getters: {
    getItem(state) {
      return (type, id) => {
        if (typeof type !== 'string') {
          throw new Error('[entities/getItem] Missing required argument.');
        }

        if (!state[type] || !id) {
          return null;
        }

        return denormalize(id, schemas[type], state);
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
        Vue.set(state, type, mergeWith({}, state[type] || {}, entities[type]), (objValue, srcValue) => {
          if (Array.isArray(objValue)) {
            return uniq(objValue.concat(srcValue));
          }

          return undefined;
        });
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
    async fetch(
      { commit },
      {
        route,
        schema,
        params,
        dataPreparer,
        isPost,
        mutationType = internalTypes.ENTITIES_UPDATE,
      },
    ) {
      let data;

      if (isPost) {
        data = await request.post(route, params);
      } else {
        data = await request.get(route, { params });
      }
      const normalizedData = normalize(dataPreparer(data), schema);
      commit(mutationType, normalizedData.entities);

      return { data, normalizedData };
    },
    async update(
      { commit },
      {
        route,
        schema,
        params,
        body,
        dataPreparer,
        mutationType = internalTypes.ENTITIES_UPDATE,
      },
    ) {
      const data = await request.put(route, body, { params });
      const normalizedData = normalize(dataPreparer(data), schema);

      commit(mutationType, normalizedData.entities);

      return { data, normalizedData };
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
  ENTITIES_DELETE: `${entitiesModuleName}/${internalTypes.ENTITIES_DELETE}`,
};

export default (store) => {
  store.registerModule(entitiesModuleName, entitiesModule);
};
