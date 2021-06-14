import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalRemediationJobMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyRemediationJobAccess() {
      return this.checkCreateAccess(USERS_PERMISSIONS.technical.remediationJob);
    },

    hasReadAnyRemediationJobAccess() {
      return this.checkReadAccess(USERS_PERMISSIONS.technical.remediationJob);
    },

    hasUpdateAnyRemediationJobAccess() {
      return this.checkUpdateAccess(USERS_PERMISSIONS.technical.remediationJob);
    },

    hasDeleteAnyRemediationJobAccess() {
      return this.checkDeleteAccess(USERS_PERMISSIONS.technical.remediationJob);
    },
  },
};
