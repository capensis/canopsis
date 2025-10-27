import { API_ROUTES } from '@/config';

import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    runAlarmFiltering() {
      return request.put(API_ROUTES.pbehavior.allPatterns);
    },

    checkPatternsPbehaviorsCount(context, { data }) {
      return request.put(API_ROUTES.pbehavior.patterns, data);
    },
  },
};
