import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsAlarmsListRemediationInstructionsFilters = {
  mixins: [authMixin],
  computed: {
    hasAccessToListRemediationInstructionsFilters() {
      return this.checkAccess(USER_PERMISSIONS.business.alarmsList.actions.listRemediationInstructionsFilters);
    },

    hasAccessToAddRemediationInstructionsFilter() {
      return this.checkAccess(USER_PERMISSIONS.business.alarmsList.actions.addRemediationInstructionsFilter);
    },

    hasAccessToEditRemediationInstructionsFilter() {
      return this.checkAccess(USER_PERMISSIONS.business.alarmsList.actions.editRemediationInstructionsFilter);
    },

    hasAccessToUserRemediationInstructionsFilter() {
      return this.checkAccess(USER_PERMISSIONS.business.alarmsList.actions.userRemediationInstructionsFilter);
    },
  },
};
