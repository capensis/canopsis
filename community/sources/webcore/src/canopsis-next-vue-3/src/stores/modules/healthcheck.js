import { defineStore } from 'pinia';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

export const useHealthcheckStore = defineStore('healthcheck', () => {
  const fetchHealthcheckStatusWithoutStore = () => request.get(API_ROUTES.healthcheck.status);

  return {
    fetchHealthcheckStatusWithoutStore,
  };
});
