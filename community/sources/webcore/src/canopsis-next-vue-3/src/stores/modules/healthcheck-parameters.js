import { defineStore } from 'pinia';

import { API_ROUTES } from '@/config';
import request from '@/services/request';

export const useHealthcheckParametersStore = defineStore('healthcheckParameters', () => {
  const fetchItemWithoutStore = () => request.get(API_ROUTES.healthcheck.parameters);
  const update = ({ data } = {}) => request.put(API_ROUTES.healthcheck.parameters, data);

  return {
    fetchItemWithoutStore,
    update,
  };
});
