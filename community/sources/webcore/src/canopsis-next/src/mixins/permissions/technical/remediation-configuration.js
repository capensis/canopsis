import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalRemediationConfigurationMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyRemediationConfigurationAccess() {
      return this.checkCreateAccess(USER_PERMISSIONS.technical.exploitation.remediationConfiguration);
    },

    hasReadAnyRemediationConfigurationAccess() {
      return this.checkReadAccess(USER_PERMISSIONS.technical.exploitation.remediationConfiguration);
    },

    hasUpdateAnyRemediationConfigurationAccess() {
      return this.checkUpdateAccess(USER_PERMISSIONS.technical.exploitation.remediationConfiguration);
    },

    hasDeleteAnyRemediationConfigurationAccess() {
      return this.checkDeleteAccess(USER_PERMISSIONS.technical.exploitation.remediationConfiguration);
    },
  },
};
