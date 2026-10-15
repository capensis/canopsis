import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsBarChartSampling = {
  mixins: [authMixin],
  computed: {
    hasAccessToSampling() {
      return this.checkAccess(USER_PERMISSIONS.business.barChart.actions.sampling);
    },
  },
};
