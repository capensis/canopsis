import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsAlarmsListFilters = {
  mixins: [authMixin],
  computed: {
    hasAccessToListFilters() {
      return this.checkAccess(USER_PERMISSIONS.business.alarmsList.actions.listFilters);
    },

    hasAccessToEditFilter() {
      return this.checkAccess(USER_PERMISSIONS.business.alarmsList.actions.editFilter);
    },

    hasAccessToAddFilter() {
      return this.checkAccess(USER_PERMISSIONS.business.alarmsList.actions.addFilter);
    },

    hasAccessToUserFilter() {
      return this.checkAccess(USER_PERMISSIONS.business.alarmsList.actions.userFilter);
    },
  },
};
