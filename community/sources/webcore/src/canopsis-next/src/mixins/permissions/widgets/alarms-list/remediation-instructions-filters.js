import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsAlarmsListRemediationInstructionsFilters = {
  mixins: [authMixin],
  computed: {
    hasAccessToRemediationInstructionsFilter() {
      return this.checkAccess(USER_PERMISSIONS.business.alarmsList.actions.remediationInstructionsFilters);
    },

    hasAccessToListRemediationInstructionsFilters() {
      return this.checkAccess(USER_PERMISSIONS.business.alarmsList.actions.listRemediationInstructionsFilters); // TODO: remove
    },

    hasAccessToAddRemediationInstructionsFilter() {
      return this.checkAccess(USER_PERMISSIONS.business.alarmsList.actions.addRemediationInstructionsFilter); // TODO: remove
    },

    hasAccessToEditRemediationInstructionsFilter() {
      return this.checkAccess(USER_PERMISSIONS.business.alarmsList.actions.editRemediationInstructionsFilter); // TODO: remove
    },

    hasAccessToUserRemediationInstructionsFilter() {
      return this.checkAccess(USER_PERMISSIONS.business.alarmsList.actions.userRemediationInstructionsFilter);
    },
  },
};
