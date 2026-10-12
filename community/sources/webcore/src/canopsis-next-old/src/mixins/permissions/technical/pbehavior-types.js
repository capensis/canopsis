import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalPbehaviorTypesMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyPbehaviorTypeAccess() {
      return this.checkCreateAccess(USER_PERMISSIONS.technical.planningType);
    },

    hasReadAnyPbehaviorTypeAccess() {
      return this.checkReadAccess(USER_PERMISSIONS.technical.planningType);
    },

    hasUpdateAnyPbehaviorTypeAccess() {
      return this.checkUpdateAccess(USER_PERMISSIONS.technical.planningType);
    },

    hasDeleteAnyPbehaviorTypeAccess() {
      return this.checkDeleteAccess(USER_PERMISSIONS.technical.planningType);
    },
  },
};
