import { API_ROUTES } from '@/config';

import { alarmDetailSchema } from '@/store/schemas';

export default {
  namespaced: true,
  actions: {
    fetchItem({ dispatch }, { data } = {}) {
      return dispatch('entities/create', {
        route: API_ROUTES.alarmDetails,
        schema: alarmDetailSchema,
        body: data,
      }, { root: true });
    },
  },
};
