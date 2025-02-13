import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalRemediationConfigurationMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyRemediationConfigurationAccess() {
      return this.checkCreateAccess(USER_PERMISSIONS.technical.remediationConfiguration);
    },

    hasReadAnyRemediationConfigurationAccess() {
      return this.checkReadAccess(USER_PERMISSIONS.technical.remediationConfiguration);
    },

    hasUpdateAnyRemediationConfigurationAccess() {
      return this.checkUpdateAccess(USER_PERMISSIONS.technical.remediationConfiguration);
    },

    hasDeleteAnyRemediationConfigurationAccess() {
      return this.checkDeleteAccess(USER_PERMISSIONS.technical.remediationConfiguration);
    },
  },
};
