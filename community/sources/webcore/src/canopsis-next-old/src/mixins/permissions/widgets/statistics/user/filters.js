import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsUserStatisticsFilters = {
  mixins: [authMixin],
  computed: {
    hasAccessToFilter() {
      return this.checkAccess(USER_PERMISSIONS.business.userStatistics.actions.filter);
    },

    hasAccessToUserFilter() {
      return this.checkAccess(USER_PERMISSIONS.business.userStatistics.actions.userFilter);
    },
  },
};
