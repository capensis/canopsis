import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsAlarmsListRemediationInstructionsFilters = {
  mixins: [authMixin],
  computed: {
    hasAccessToRemediationInstructionsFilter() {
      return this.checkAccess(USER_PERMISSIONS.business.alarmsList.actions.remediationInstructionsFilters);
    },

    hasAccessToUserRemediationInstructionsFilter() {
      return this.checkAccess(USER_PERMISSIONS.business.alarmsList.actions.userRemediationInstructionsFilter);
    },
  },
};
