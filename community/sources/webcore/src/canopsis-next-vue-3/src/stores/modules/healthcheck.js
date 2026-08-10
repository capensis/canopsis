import { defineStore } from 'pinia';

import { API_ROUTES } from '@/config';
import request from '@/services/request';

export const useHealthcheckStore = defineStore('healthcheck', () => {
  const fetchHealthcheckEnginesWithoutStore = () => request.get(API_ROUTES.healthcheck.engines);
  const fetchHealthcheckStatusWithoutStore = () => request.get(API_ROUTES.healthcheck.status);
  const fetchHealthcheckEnginesOrderWithoutStore = () =>
    request.get(API_ROUTES.healthcheck.enginesOrder);

  return {
    fetchHealthcheckEnginesWithoutStore,
    fetchHealthcheckStatusWithoutStore,
    fetchHealthcheckEnginesOrderWithoutStore,
  };
});
