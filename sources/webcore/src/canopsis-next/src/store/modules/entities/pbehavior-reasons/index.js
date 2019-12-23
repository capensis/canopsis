import { API_ROUTES } from '@/config';

import request from '@/services/request';

export const types = {
  FETCH_PBEHAVIOR_REASONS: 'FETCH_PBEHAVIOR_REASONS',
  FETCH_PBEHAVIOR_REASONS_COMPLETED: 'FETCH_PBEHAVIOR_REASONS_COMPLETED',
};

export default {
  namespaced: true,
  state: {
    pending: false,
    pbehaviorReasons: [],
  },
  getters: {
    pending: state => state.pending,
    pbehaviorReasons: state => state.pbehaviorReasons,
  },
  mutations: {
    [types.FETCH_PBEHAVIOR_REASONS](state) {
      state.pending = true;
    },
    [types.FETCH_PBEHAVIOR_REASONS_COMPLETED](state, { reasons }) {
      state.pending = false;
      state.pbehaviorReasons = reasons;
    },
  },
  actions: {
    async fetchPbehaviorReasons({ commit }) {
      try {
        commit(types.FETCH_PBEHAVIOR_REASONS);

        const { reasons } = await request.get(API_ROUTES.pbehaviorReasons);

        commit(types.FETCH_PBEHAVIOR_REASONS_COMPLETED, { reasons });
      } catch (err) {
        console.error(err);
      }
    },
  },
};
