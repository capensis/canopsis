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
          return null;
        }

        return denormalize(ids, [schemas[type]], state.byId);
      };
    },
  },
  mutations: {
    [types.ENTITIES_UPDATE](state, entities) {
      Object.keys(entities).forEach((key) => {
        Vue.set(state.byId, key, {
          ...(state.byId[key] || {}),
          ...entities[key],
        });
      });
    },
    [types.ENTITIES_MERGE](state, entities) {
      Object.keys(entities).forEach((key) => {
        Vue.set(state.byId, key, merge({}, state.byId[key] || {}, entities[key]));
      });
    },
    [types.ENTITIES_DELETE](state, entities) {
      Object.keys(entities).forEach((key) => {
        Vue.set(state.byId, key, omit(state.byId[key], Object.keys(entities[key])));
      });
    },
  },
};
