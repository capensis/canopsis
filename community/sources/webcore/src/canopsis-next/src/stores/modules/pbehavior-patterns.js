import { defineStore } from 'pinia';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

export const usePbehaviorPatternsStore = defineStore('pbehavior-patterns', () => {
  const runAlarmFiltering = () => request.put(API_ROUTES.pbehavior.allPatterns);

  const checkPatternsPbehaviorsCount = ({ data } = {}) =>
    request.put(API_ROUTES.pbehavior.patterns, data);

  return {
    runAlarmFiltering,
    checkPatternsPbehaviorsCount,
  };
});
