import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsPieChartSampling = {
  mixins: [authMixin],
  computed: {
    hasAccessToSampling() {
      return this.checkAccess(USERS_PERMISSIONS.business.pieChart.actions.sampling);
    },
  },
};
