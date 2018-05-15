import Vue from 'vue';
import omit from 'lodash/omit';
import merge from 'lodash/merge';
import { denormalize } from 'normalizr';

import schemas from '@/store/schemas';

import types from './types';
import alarmModule from './alarm';
import pbehaviorModule from './pbehavior';

export default {
  namespaced: true,
  state: {
    byId: {},
  },
  modules: {
    pbehavior: pbehaviorModule,
    alarm: alarmModule,
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
    [types.ENTITIES_UPDATE](state, entities) {
      Object.keys(entities).forEach((type) => {
        Vue.set(state.byId, type, {
          ...(state.byId[type] || {}),
          ...entities[type],
        });
      });
    },
    [types.ENTITIES_MERGE](state, entities) {
      Object.keys(entities).forEach((type) => {
        Vue.set(state.byId, type, merge({}, state.byId[type] || {}, entities[type]));
      });
    },
    [types.ENTITIES_DELETE](state, entities) {
      Object.keys(entities).forEach((type) => {
        entities[type].forEach((id) => {
          const entity = state.byId[type][id];

          if (entity && entity._embedded) {
            const { parentType, parentId, relationType } = entity._embedded;

            if (state.byId[parentType]) {
              const parentEntity = state.byId[parentType][parentId];

              if (parentEntity) {
                parentEntity[relationType] = parentEntity[relationType].filter(v => v !== id);
              }
            }
          }
        });

        Vue.set(state.byId, type, omit(state.byId[type], entities[type]));
      });
    },
  },
};
