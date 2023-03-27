import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsLineChartInterval = {
  mixins: [authMixin],
  computed: {
    hasAccessToInterval() {
      return this.checkAccess(USERS_PERMISSIONS.business.lineChart.actions.interval);
    },
  },
};
