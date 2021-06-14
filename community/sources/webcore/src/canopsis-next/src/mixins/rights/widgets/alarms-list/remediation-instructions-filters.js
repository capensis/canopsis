import { USERS_RIGHTS } from '@/constants';

import authMixin from '@/mixins/auth';

export default {
  mixins: [authMixin],
  computed: {
    hasAccessToListRemediationInstructionsFilters() {
      return this.checkAccess(USERS_RIGHTS.business.alarmsList.actions.listRemediationInstructionsFilters);
    },

    hasAccessToAddRemediationInstructionsFilter() {
      return this.checkAccess(USERS_RIGHTS.business.alarmsList.actions.addRemediationInstructionsFilter);
    },

    hasAccessToEditRemediationInstructionsFilter() {
      return this.checkAccess(USERS_RIGHTS.business.alarmsList.actions.editRemediationInstructionsFilter);
    },

    hasAccessToUserRemediationInstructionsFilter() {
      return this.checkAccess(USERS_RIGHTS.business.alarmsList.actions.userRemediationInstructionsFilter);
    },
  },
};
