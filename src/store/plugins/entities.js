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
     * @param {Object.<string, Array.<string>>} entitiesIds - Object of entities ids
     */
    [internalTypes.ENTITIES_DELETE](state, entitiesIds) {
      Object.keys(entitiesIds).forEach((type) => {
        entitiesIds[type].forEach((id) => {
          const entity = state[type][id];
          const { parentType, parentId, relationType } = get(entity, '_embedded', {});
          const parentEntity = get(state, [parentType, parentId]);

          if (parentEntity) {
            Vue.set(parentEntity, relationType, parentEntity[relationType].filter(v => v !== id));
          }
        });

        Vue.set(state, type, omit(state[type], entitiesIds[type]));
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
        mutationType = internalTypes.ENTITIES_UPDATE,
      },
    ) {
      const [data] = await request.get(route, { params });
      const normalizedData = normalize(dataPreparer(data), schema);

      commit(mutationType, normalizedData.entities);

      return { data, normalizedData };
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
