import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalPbehaviorReasonsMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyPbehaviorReasonAccess() {
      return this.checkCreateAccess(USER_PERMISSIONS.technical.planningReason);
    },

    hasReadAnyPbehaviorReasonAccess() {
      return this.checkReadAccess(USER_PERMISSIONS.technical.planningReason);
    },

    hasUpdateAnyPbehaviorReasonAccess() {
      return this.checkUpdateAccess(USER_PERMISSIONS.technical.planningReason);
    },

    hasDeleteAnyPbehaviorReasonAccess() {
      return this.checkDeleteAccess(USER_PERMISSIONS.technical.planningReason);
    },
  },
};
