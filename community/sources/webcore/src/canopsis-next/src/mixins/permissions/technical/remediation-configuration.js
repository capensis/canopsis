import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalRemediationConfigurationMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyRemediationConfigurationAccess() {
      return this.checkCreateAccess(USERS_PERMISSIONS.technical.remediationConfiguration);
    },

    hasReadAnyRemediationConfigurationAccess() {
      return this.checkReadAccess(USERS_PERMISSIONS.technical.remediationConfiguration);
    },

    hasUpdateAnyRemediationConfigurationAccess() {
      return this.checkUpdateAccess(USERS_PERMISSIONS.technical.remediationConfiguration);
    },

    hasDeleteAnyRemediationConfigurationAccess() {
      return this.checkDeleteAccess(USERS_PERMISSIONS.technical.remediationConfiguration);
    },
  },
};
