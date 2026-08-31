import { defineStore } from 'pinia';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

export const useTechMetricsStore = defineStore('techMetrics', () => {
  const createTechMetricsExport = () => request.post(API_ROUTES.techMetrics);
  const fetchTechMetricsExport = () => request.get(API_ROUTES.techMetrics);
  const fetchTechMetricsSettings = () => request.get(API_ROUTES.techMetricsSettings);
  const updateTechMetricsSettings = ({ data } = {}) => request.put(API_ROUTES.techMetricsSettings, data);

  return {
    createTechMetricsExport,
    fetchTechMetricsExport,
    fetchTechMetricsSettings,
    updateTechMetricsSettings,
  };
});
