import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalRemediationInstructionMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyRemediationInstructionAccess() {
      return this.checkCreateAccess(USER_PERMISSIONS.technical.remediationInstruction);
    },

    hasReadAnyRemediationInstructionAccess() {
      return this.checkReadAccess(USER_PERMISSIONS.technical.remediationInstruction);
    },

    hasUpdateAnyRemediationInstructionAccess() {
      return this.checkUpdateAccess(USER_PERMISSIONS.technical.remediationInstruction);
    },

    hasDeleteAnyRemediationInstructionAccess() {
      return this.checkDeleteAccess(USER_PERMISSIONS.technical.remediationInstruction);
    },
  },
};
