import { normalize } from 'normalizr';

import request from '@/services/request';
import i18n from '@/i18n';
import schemas from '@/store/schemas';
import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';
import { types as entitiesTypes } from '@/store/plugins/entities';

const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
  FETCH_BY_ID: 'FETCH_BY_ID',
  FETCH_BY_ID_COMPLETED: 'FETCH_BY_ID_COMPLETED',
  FETCH_BY_ID_FAILED: 'FETCH_BY_ID_FAILED',
};

export default {
  namespaced: true,
  state: {
    allIds: [],
    error: '',
    pending: false,
    meta: {},
  },
  getters: {
    allIds: state => state.allIds,
    items: (state, getters, rootState, rootGetters) =>
      rootGetters['entities/getList'](ENTITIES_TYPES.pbehavior, state.allIds),
    error: state => state.error,
    pending: state => state.pending,
    meta: state => state.meta,
  },
  mutations: {
    [types.FETCH_LIST](state) {
      state.pending = true;
    },
    [types.FETCH_LIST_COMPLETED](state, { allIds, meta }) {
      state.allIds = allIds;
      state.meta = meta;
      state.pending = false;
    },
    [types.FETCH_LIST_FAILED](state) {
      state.pending = false;
    },
    [types.FETCH_BY_ID](state) {
      state.pending = true;
    },
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
    async fetchList({ dispatch, commit }, { params } = {}) {
      try {
        commit(types.FETCH_LIST);

        const { data, normalizedData } = await dispatch('entities/fetch', {
          route: API_ROUTES.pbehavior.list,
          schema: [schemas.pbehavior],
          params,
          dataPreparer: d => d.data,
        }, { root: true });

        commit(types.FETCH_LIST_COMPLETED, {
          allIds: normalizedData.result,
          meta: {
            total: data.total,
          },
        });
      } catch (err) {
        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });
        commit(types.FETCH_LIST_FAILED);
      }
    },

    async fetchListByEntityId({ commit, dispatch }, { id }) {
      try {
        const { normalizedData } = await dispatch('entities/fetch', {
          route: `${API_ROUTES.pbehaviorById}/${id}`,
          schema: [schemas.pbehavior],
        }, { root: true });
        commit(types.FETCH_BY_ID_COMPLETED, normalizedData.result);
      } catch (err) {
        commit(types.FETCH_BY_ID_FAILED, err);
        console.warn(err);
      }
    },

    async create({ commit }, { data, parents, parentsType }) {
      try {
        const id = await request.post(API_ROUTES.pbehavior.pbehavior, data);
        const pbehavior = {
          ...data,
          enabled: true,
          _id: id,
        };

        if (parents && parentsType) {
          const parentSchema = schemas[parentsType];

          const parentEntities = parents
            .map(parent => ({
              ...parent,
              pbehaviors: parent.pbehaviors ? [...parent.pbehaviors, pbehavior] : [pbehavior],
            }));

          const { entities } = normalize(parentEntities, [parentSchema]);

          commit(entitiesTypes.ENTITIES_MERGE, entities, { root: true });
        }
      } catch (err) {
        console.error(err);

        throw err;
      }
    },

    async remove({ dispatch }, { id }) {
      try {
        await request.delete(`${API_ROUTES.pbehavior.pbehavior}/${id}`);
        await dispatch('entities/removeFromStore', {
          id,
          type: ENTITIES_TYPES.pbehavior,
        }, { root: true });
      } catch (err) {
        console.warn(err);
      }
    },
  },
};

