import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsAlarmStatisticsFilters = {
  mixins: [authMixin],
  computed: {
    hasAccessToListFilters() {
      return this.checkAccess(USERS_PERMISSIONS.business.alarmStatistics.actions.listFilters);
    },

    hasAccessToEditFilter() {
      return this.checkAccess(USERS_PERMISSIONS.business.alarmStatistics.actions.editFilter);
    },

    hasAccessToAddFilter() {
      return this.checkAccess(USERS_PERMISSIONS.business.alarmStatistics.actions.addFilter);
    },

    hasAccessToUserFilter() {
      return this.checkAccess(USERS_PERMISSIONS.business.alarmStatistics.actions.userFilter);
    },
  },
};
