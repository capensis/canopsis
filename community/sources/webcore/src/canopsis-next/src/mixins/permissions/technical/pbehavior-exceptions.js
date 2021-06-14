import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalPbehaviorExceptionsMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyPbehaviorExceptionAccess() {
      return this.checkCreateAccess(USERS_PERMISSIONS.technical.planningExceptions);
    },

    hasReadAnyPbehaviorExceptionAccess() {
      return this.checkReadAccess(USERS_PERMISSIONS.technical.planningExceptions);
    },

    hasUpdateAnyPbehaviorExceptionAccess() {
      return this.checkUpdateAccess(USERS_PERMISSIONS.technical.planningExceptions);
    },

    hasDeleteAnyPbehaviorExceptionAccess() {
      return this.checkDeleteAccess(USERS_PERMISSIONS.technical.planningExceptions);
    },
  },
};
