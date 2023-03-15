import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsLineChartSampling = {
  mixins: [authMixin],
  computed: {
    hasAccessToSampling() {
      return this.checkAccess(USERS_PERMISSIONS.business.lineChart.actions.sampling);
    },
  },
};
