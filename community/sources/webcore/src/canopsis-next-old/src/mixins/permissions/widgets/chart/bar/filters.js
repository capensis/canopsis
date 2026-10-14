import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsBarChartFilters = {
  mixins: [authMixin],
  computed: {
    hasAccessToFilter() {
      return this.checkAccess(USER_PERMISSIONS.business.barChart.actions.filter);
    },

    hasAccessToUserFilter() {
      return this.checkAccess(USER_PERMISSIONS.business.barChart.actions.userFilter);
    },
  },
};
