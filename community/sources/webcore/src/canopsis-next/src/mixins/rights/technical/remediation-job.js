import { USERS_RIGHTS } from '@/constants';

import authMixin from '@/mixins/auth';

export default {
  mixins: [authMixin],
  computed: {
    hasCreateAnyRemediationJobAccess() {
      return this.checkCreateAccess(USERS_RIGHTS.technical.remediationJob);
    },

    hasReadAnyRemediationJobAccess() {
      return this.checkReadAccess(USERS_RIGHTS.technical.remediationJob);
    },

    hasUpdateAnyRemediationJobAccess() {
      return this.checkUpdateAccess(USERS_RIGHTS.technical.remediationJob);
    },

    hasDeleteAnyRemediationJobAccess() {
      return this.checkDeleteAccess(USERS_RIGHTS.technical.remediationJob);
    },
  },
};
