import request from '@/services/request';
import { API_ROUTES } from '@/config';

import { types as entitiesTypes } from '@/store/plugins/entities';

const types = {
  FETCH_BY_ID_COMPLETED: 'FETCH_BY_ID_COMPLETED',
  FETCH_BY_ID_FAILED: 'FETCH_BY_ID_FAILED',
};

export default {
  namespaced: true,
  state: {
    pbehaviorsList: [],
    error: '',
    pending: true,
  },
  getters: {
    pbehaviorsList: state => state.pbehaviorsList,
    error: state => state.error,
    pending: state => state.pending,
  },
  mutations: {
    [types.FETCH_BY_ID_COMPLETED](state, payload) {
      state.pbehaviorsList = payload;
      state.pending = false;
    },
    [types.FETCH_BY_ID_FAILED](state, err) {
      state.error = err;
      state.pending = false;
    },
  },
  actions: {
    async create(context, data) {
      try {
        await request.post(API_ROUTES.pbehavior, data);
      } catch (err) {
        console.error(err);

        throw err;
      }
    },
    async remove({ commit }, { id }) {
      try {
        await request.delete(`${API_ROUTES.pbehavior}/${id}`);

        commit(
          entitiesTypes.ENTITIES_DELETE,
          { pbehavior: [id] },
          { root: true },
        );
      } catch (err) {
        console.warn(err);
      }
    },
    async fetchById({ commit }, { id }) {
      try {
        const data = await request.get(`${API_ROUTES.pbehaviorById}/${id}`);
        commit(types.FETCH_BY_ID_COMPLETED, data);
      } catch (err) {
        commit(types.FETCH_BY_ID_FAILED, err);
        console.warn(err);
      }
    },
  },
};
