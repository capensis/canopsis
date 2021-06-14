import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalPbehaviorReasonsMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyPbehaviorReasonAccess() {
      return this.checkCreateAccess(USERS_PERMISSIONS.technical.planningReason);
    },

    hasReadAnyPbehaviorReasonAccess() {
      return this.checkReadAccess(USERS_PERMISSIONS.technical.planningReason);
    },

    hasUpdateAnyPbehaviorReasonAccess() {
      return this.checkUpdateAccess(USERS_PERMISSIONS.technical.planningReason);
    },

    hasDeleteAnyPbehaviorReasonAccess() {
      return this.checkDeleteAccess(USERS_PERMISSIONS.technical.planningReason);
    },
  },
};
