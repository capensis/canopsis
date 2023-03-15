import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsBarChartInterval = {
  mixins: [authMixin],
  computed: {
    hasAccessToInterval() {
      return this.checkAccess(USERS_PERMISSIONS.business.barChart.actions.interval);
    },
  },
};
