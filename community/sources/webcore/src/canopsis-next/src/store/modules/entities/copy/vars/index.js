import { API_ROUTES } from '@/config';

import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    fetchEventFiltersVarsWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.copyVars.eventFilters, { params });
    },

    fetchDynamicInfosVarsWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.copyVarsCat.dynamicInfos, { params });
    },
  },
};
