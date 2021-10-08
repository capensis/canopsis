import { API_ROUTES } from '@/config';

import request from '@/services/request';
import localStorageService from '@/services/local-storage';

export const types = {
  SET_VIEW_STATS_ID: 'SET_VIEW_STATS_ID',
};

const LOCAL_STORAGE_KEY = 'viewStatsId';

export default {
  namespaced: true,
  state: {
    currentViewStatsId: localStorageService.get(LOCAL_STORAGE_KEY),
  },
  mutations: {
    [types.SET_VIEW_STATS_ID]: (state, id) => {
      state.currentViewStatsId = id;
    },
  },
  actions: {
    async create({ commit }) {
      const { _id: id } = await request.post(API_ROUTES.viewStats);

      commit(types.SET_VIEW_STATS_ID, id);
      localStorageService.set(LOCAL_STORAGE_KEY, id);
    },

    async update({ state, rootGetters }, { data }) {
      if (!state.currentViewStatsId || !rootGetters['auth/isLoggedIn']) {
        return Promise.resolve();
      }

      return request.put(`${API_ROUTES.viewStats}/${state.currentViewStatsId}`, data);
    },
  },
};
