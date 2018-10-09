import { normalize } from 'normalizr';

import request from '@/services/request';
import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';
import schemas from '@/store/schemas';
import { types as entitiesTypes } from '@/store/plugins/entities';

const types = {
  FETCH_BY_ID_COMPLETED: 'FETCH_BY_ID_COMPLETED',
  FETCH_BY_ID_FAILED: 'FETCH_BY_ID_FAILED',
};

export default {
  namespaced: true,
  state: {
    allIds: [],
    error: '',
    pending: true,
  },
  getters: {
    allIds: state => state.allIds,
    items: (state, getters, rootState, rootGetters) => rootGetters['entities/getList'](ENTITIES_TYPES.pbehavior, state.allIds),
    error: state => state.error,
    pending: state => state.pending,
  },
  mutations: {
    [types.FETCH_BY_ID_COMPLETED](state, ids) {
      state.allIds = ids;
      state.pending = false;
    },
    [types.FETCH_BY_ID_FAILED](state, err) {
      state.error = err;
      state.pending = false;
    },
  },
  actions: {
    async create({ commit }, { data, parents, parentsType }) {
      try {
        const parentSchema = schemas[parentsType];
        const id = await request.post(API_ROUTES.pbehavior, data);
        const pbehavior = {
          ...data,
          enabled: true,
          _id: id,
        };

        const parentEntities = parents.map(parent => ({ ...parent, pbehaviors: [...parent.pbehaviors, pbehavior] }));

        const { entities } = normalize(parentEntities, [parentSchema]);

        commit(entitiesTypes.ENTITIES_MERGE, entities, { root: true });
      } catch (err) {
        console.error(err);

        throw err;
      }
    },
    async remove({ dispatch }, { id }) {
      try {
        await request.delete(`${API_ROUTES.pbehavior}/${id}`);
        await dispatch('entities/removeFromStore', {
          id,
          type: ENTITIES_TYPES.pbehavior,
        }, { root: true });
      } catch (err) {
        console.warn(err);
      }
    },
    fetch({ dispatch }, { id }) {
      return dispatch('entities/fetch', {
        route: `${API_ROUTES.pbehaviorById}/${id}`,
        schema: [schemas.pbehavior],
      }, { root: true });
    },
    async fetchListByEntityId({ commit, dispatch }, { id }) {
      try {
        const { normalizedData } = await dispatch('fetch', { id });
        commit(types.FETCH_BY_ID_COMPLETED, normalizedData.result);
      } catch (err) {
        commit(types.FETCH_BY_ID_FAILED, err);
        console.warn(err);
      }
    },
  },
};

