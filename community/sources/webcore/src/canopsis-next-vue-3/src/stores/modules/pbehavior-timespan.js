import { defineStore } from 'pinia';

import { API_ROUTES } from '@/config';
import request from '@/services/request';

export const usePbehaviorTimespanStore = defineStore('pbehaviorTimespan', () => {
  const fetchListWithoutStore = ({ data } = {}) =>
    request.post(API_ROUTES.pbehavior.timespan, data);

  return {
    fetchListWithoutStore,
  };
});
