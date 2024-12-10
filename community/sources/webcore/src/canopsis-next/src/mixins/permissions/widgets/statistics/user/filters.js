import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsUserStatisticsFilters = {
  mixins: [authMixin],
  computed: {
    hasAccessToListFilters() {
      return this.checkAccess(USER_PERMISSIONS.business.userStatistics.actions.listFilters);
    },

    hasAccessToEditFilter() {
      return this.checkAccess(USER_PERMISSIONS.business.userStatistics.actions.editFilter);
    },

    hasAccessToAddFilter() {
      return this.checkAccess(USER_PERMISSIONS.business.userStatistics.actions.addFilter);
    },

    hasAccessToUserFilter() {
      return this.checkAccess(USER_PERMISSIONS.business.userStatistics.actions.userFilter);
    },
  },
};
