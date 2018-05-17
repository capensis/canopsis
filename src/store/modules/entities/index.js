import Vue from 'vue';
import omit from 'lodash/omit';
import merge from 'lodash/merge';
import get from 'lodash/get';
import { denormalize } from 'normalizr';

import schemas from '@/store/schemas';

import types from './types';
import alarmModule from './alarm';
import pbehaviorModule from './pbehavior';

export default {
  namespaced: true,
  modules: {
    pbehavior: pbehaviorModule,
    alarm: alarmModule,
  },
  state: {
    byId: {},
  },
  getters: {
    getItem(state) {
      return (type, id) => {
        if (typeof type !== 'string' || !id) {
          throw new Error('[entities/getItem] Missing required argument.');
        }

        if (!state[type]) {
          return null;
        }

        return denormalize(id, schemas[type], state.byId);
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

        const result = denormalize(ids, [schemas[type]], state.byId);

        return result.filter(v => !!v);
      };
    },
  },
  mutations: {
    /**
     * @param {Object} state - state of the module
     * @param {Object.<string, Object>} entities - Object of entities
     */
    [types.ENTITIES_UPDATE](state, entities) {
      Object.keys(entities).forEach((type) => {
        Vue.set(state.byId, type, {
          ...(state.byId[type] || {}),
          ...entities[type],
        });
      });
    },

    /**
     * @param {Object} state - state of the module
     * @param {Object.<string, Object>} entities - Object of entities
     */
    [types.ENTITIES_MERGE](state, entities) {
      Object.keys(entities).forEach((type) => {
        Vue.set(state.byId, type, merge({}, state.byId[type] || {}, entities[type]));
      });
    },

    /**
     * @param {Object} state - state of the module
     * @param {Object.<string, Array.<string>>} entitiesIds - Object of entities ids
     */
    [types.ENTITIES_DELETE](state, entitiesIds) {
      Object.keys(entitiesIds).forEach((type) => {
        entitiesIds[type].forEach((id) => {
          const entity = state.byId[type][id];
          const { parentType, parentId, relationType } = get(entity, '_embedded', {});
          const parentEntity = get(state.byId, [parentType, parentId]);

          if (parentEntity) {
            Vue.set(parentEntity, relationType, parentEntity[relationType].filter(v => v !== id));
          }
        });

        Vue.set(state.byId, type, omit(state.byId[type], entitiesIds[type]));
      });
    },
  },
};
