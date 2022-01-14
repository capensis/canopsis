import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import { createEntityModule } from '@/store/plugins/entities';

import request from '@/services/request';

export const types = {
  SET_UPDATED_AT: 'SET_UPDATED_AT',
};

export default createEntityModule({
  route: API_ROUTES.ratingSettings,
  entityType: ENTITIES_TYPES.ratingSettings,
  dataPreparer: d => d.data,
  withFetchingParams: true,
  withMeta: true,
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
    fetchListWithoutStore(context, { params }) {
      return request.get(API_ROUTES.ratingSettings, { params });
    },

    async bulkUpdate({ commit }, { data }) {
      await request.put(API_ROUTES.bulkRatingSettings, data);

      commit(types.SET_UPDATED_AT);
    },
  },
});
