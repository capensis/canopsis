import { defineStore } from 'pinia';

import { API_ROUTES } from '@/config';
import request from '@/services/request';

export const useCopyVarsStore = defineStore('copyVars', () => {
  const fetchEventFiltersCopyVarsWithoutStore = ({ params } = {}) =>
    request.get(API_ROUTES.copyVars.eventFilters, { params });

  const fetchDynamicInfosCopyVarsWithoutStore = ({ params } = {}) =>
    request.get(API_ROUTES.copyVarsCat.dynamicInfos, { params });

  return {
    fetchEventFiltersCopyVarsWithoutStore,
    fetchDynamicInfosCopyVarsWithoutStore,
  };
});
