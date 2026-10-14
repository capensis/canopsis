import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsLineChartSampling = {
  mixins: [authMixin],
  computed: {
    hasAccessToSampling() {
      return this.checkAccess(USER_PERMISSIONS.business.lineChart.actions.sampling);
    },
  },
};
