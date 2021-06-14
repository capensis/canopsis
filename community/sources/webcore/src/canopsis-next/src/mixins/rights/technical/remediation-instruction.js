import { USERS_RIGHTS } from '@/constants';

import authMixin from '@/mixins/auth';

export default {
  mixins: [authMixin],
  computed: {
    hasCreateAnyRemediationInstructionAccess() {
      return this.checkCreateAccess(USERS_RIGHTS.technical.remediationInstruction);
    },

    hasReadAnyRemediationInstructionAccess() {
      return this.checkReadAccess(USERS_RIGHTS.technical.remediationInstruction);
    },

    hasUpdateAnyRemediationInstructionAccess() {
      return this.checkUpdateAccess(USERS_RIGHTS.technical.remediationInstruction);
    },

    hasDeleteAnyRemediationInstructionAccess() {
      return this.checkDeleteAccess(USERS_RIGHTS.technical.remediationInstruction);
    },
  },
};
