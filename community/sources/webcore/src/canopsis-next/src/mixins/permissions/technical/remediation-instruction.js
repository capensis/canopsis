import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalRemediationInstructionMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyRemediationInstructionAccess() {
      return this.checkCreateAccess(USERS_PERMISSIONS.technical.remediationInstruction);
    },

    hasReadAnyRemediationInstructionAccess() {
      return this.checkReadAccess(USERS_PERMISSIONS.technical.remediationInstruction);
    },

    hasUpdateAnyRemediationInstructionAccess() {
      return this.checkUpdateAccess(USERS_PERMISSIONS.technical.remediationInstruction);
    },

    hasDeleteAnyRemediationInstructionAccess() {
      return this.checkDeleteAccess(USERS_PERMISSIONS.technical.remediationInstruction);
    },
  },
};
