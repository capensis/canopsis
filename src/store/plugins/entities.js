import Vue from 'vue';
import omit from 'lodash/omit';
import merge from 'lodash/merge';
import get from 'lodash/get';
import { normalize, denormalize } from 'normalizr';

import request from '@/services/request';
import schemas from '@/store/schemas';

const entitiesModuleName = 'entities';

const internalTypes = {
  ENTITIES_UPDATE: 'ENTITIES_UPDATE',
  ENTITIES_MERGE: 'ENTITIES_MERGE',
  ENTITIES_DELETE: 'ENTITIES_DELETE',
};

// TODO: move and finish it
const prepareDataForDelete = (schema, data, filter = () => {}) => {
  let entitiesToMerge = {};
  let entitiesToDelete = {
    [schema.key]: {
      [data._id]: {},
    },
  };

  Object.keys(schema.schema).forEach((key) => {
    if (Array.isArray(schema.schema[key])) {
      const childrenSchema = schema.schema[key][0];

      data[key].forEach((entity) => {
        const parents = get(entity, '_embedded.parents', []);

        if (parents.length <= 1) {
          const result = prepareDataForDelete(
            childrenSchema,
            entity,
            v => v.id !== entity._id || (v.id === entity._id && v.type !== childrenSchema.key),
          );

          entitiesToMerge = {
            ...entitiesToMerge,
            ...result.entitiesToMerge,
          };

          entitiesToDelete = {
            ...entitiesToDelete,
            ...result.entitiesToDelete,
          };
        } else {
          if (!entitiesToMerge[childrenSchema.key]) {
            entitiesToMerge[childrenSchema.key] = {};
          }

          entitiesToMerge[childrenSchema.key][entity[childrenSchema.idAttribute]] = {
            ...entity,
            _embedded: {
              ...entity._embedded,
              parents: parents.filter(filter),
            },
          };
        }
      });
    } else {
      // TODO: finish this method
    }
  });

  return {
    entitiesToMerge,
    entitiesToDelete,
  };
};

const entitiesModule = {
  namespaced: true,
  getters: {
    getItem(state) {
      return (type, id) => {
        if (typeof type !== 'string' || !id) {
          throw new Error('[entities/getItem] Missing required argument.');
        }

        if (!state[type]) {
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
        Vue.set(state, type, merge({}, state[type] || {}, entities[type]));
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
    remove({ commit, getters, state }, { id, type }) {
      const schema = schemas[type];
      const item = getters.getItem(type, id);
      const parents = get(item, '_embedded.parents', []);

      const {
        entitiesToMerge,
        entitiesToDelete,
      } = prepareDataForDelete(schema, item, v => v.id !== id || (v.id === id && v.type !== type));

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
