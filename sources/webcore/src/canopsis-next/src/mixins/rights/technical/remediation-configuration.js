import { USERS_RIGHTS } from '@/constants';

import authMixin from '@/mixins/auth';

export default {
  mixins: [authMixin],
  computed: {
    hasCreateAnyRemediationConfigurationAccess() {
      return this.checkCreateAccess(USERS_RIGHTS.technical.remediationConfiguration);
    },

    hasReadAnyRemediationConfigurationAccess() {
      return this.checkReadAccess(USERS_RIGHTS.technical.remediationConfiguration);
    },

    hasUpdateAnyRemediationConfigurationAccess() {
      return this.checkUpdateAccess(USERS_RIGHTS.technical.remediationConfiguration);
    },

    hasDeleteAnyRemediationConfigurationAccess() {
      return this.checkDeleteAccess(USERS_RIGHTS.technical.remediationConfiguration);
    },
  },
};
