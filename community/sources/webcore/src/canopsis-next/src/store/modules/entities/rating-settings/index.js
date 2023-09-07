import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

export const types = {
  SET_UPDATED_AT: 'SET_UPDATED_AT',
};

export default createCRUDModule({
  route: API_ROUTES.ratingSettings,
  withFetchingParams: true,
  withWithoutStore: true,
}, {
  getters: {
    updatedAt: state => state.updatedAt,
  },
  state: {
    updatedAt: null,
  },
  mutations: {
    [types.SET_UPDATED_AT](state) {
      state.updatedAt = Date.now();
    },
  },
  actions: {
    async bulkUpdate({ commit }, { data }) {
      await request.put(API_ROUTES.bulkRatingSettings, data);

      commit(types.SET_UPDATED_AT);
    },
  },
});
