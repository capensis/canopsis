import { API_ROUTES } from '@/config';

import { alarmDetailSchema } from '@/store/schemas';
import { ENTITIES_TYPES } from '@/constants';

export default {
  namespaced: true,
  state: {
    queries: {},
  },
  getters: {
    getItem: (state, getters, rootState, rootGetters) => id => rootGetters['entities/getItem'](
      ENTITIES_TYPES.alarmDetail,
      id,
    ),
  },
  actions: {
    fetchItem({ dispatch }, { data } = {}) {
      return dispatch('entities/create', {
        route: API_ROUTES.alarmDetails,
        schema: [alarmDetailSchema],
        body: data,
      }, { root: true });
    },
  },
};
