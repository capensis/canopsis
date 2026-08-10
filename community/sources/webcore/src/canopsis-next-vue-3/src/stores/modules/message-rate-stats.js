import { defineStore } from 'pinia';

import { API_ROUTES } from '@/config';
import request from '@/services/request';

export const useMessageRateStatsStore = defineStore('messageRateStats', () => {
  const fetchMessageRateStatsWithoutStore = ({ params } = {}) =>
    request.get(API_ROUTES.messageRateStats, {
      params,
    });

  return {
    fetchMessageRateStatsWithoutStore,
  };
});
