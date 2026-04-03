import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalRemediationJobMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyRemediationJobAccess() {
      return this.checkCreateAccess(USER_PERMISSIONS.technical.exploitation.remediationJob);
    },

    hasReadAnyRemediationJobAccess() {
      return this.checkReadAccess(USER_PERMISSIONS.technical.exploitation.remediationJob);
    },

    hasUpdateAnyRemediationJobAccess() {
      return this.checkUpdateAccess(USER_PERMISSIONS.technical.exploitation.remediationJob);
    },

    hasDeleteAnyRemediationJobAccess() {
      return this.checkDeleteAccess(USER_PERMISSIONS.technical.exploitation.remediationJob);
    },
  },
};
