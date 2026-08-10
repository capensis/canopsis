import { defineStore } from 'pinia';

import { API_ROUTES } from '@/config';
import request from '@/services/request';

export const useMetricsSettingsStore = defineStore('metricsSettings', () => {
  const fetchItemWithoutStore = () => request.get(API_ROUTES.metrics.settings);
  const update = ({ data } = {}) => request.put(API_ROUTES.metrics.settings, data);

  return {
    fetchItemWithoutStore,
    update,
  };
});
