import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalPbehaviorTypesMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyPbehaviorTypeAccess() {
      return this.checkCreateAccess(USERS_PERMISSIONS.technical.planningType);
    },

    hasReadAnyPbehaviorTypeAccess() {
      return this.checkReadAccess(USERS_PERMISSIONS.technical.planningType);
    },

    hasUpdateAnyPbehaviorTypeAccess() {
      return this.checkUpdateAccess(USERS_PERMISSIONS.technical.planningType);
    },

    hasDeleteAnyPbehaviorTypeAccess() {
      return this.checkDeleteAccess(USERS_PERMISSIONS.technical.planningType);
    },
  },
};
