import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalPbehaviorExceptionsMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyPbehaviorExceptionAccess() {
      return this.checkCreateAccess(USER_PERMISSIONS.technical.planningExceptions);
    },

    hasReadAnyPbehaviorExceptionAccess() {
      return this.checkReadAccess(USER_PERMISSIONS.technical.planningExceptions);
    },

    hasUpdateAnyPbehaviorExceptionAccess() {
      return this.checkUpdateAccess(USER_PERMISSIONS.technical.planningExceptions);
    },

    hasDeleteAnyPbehaviorExceptionAccess() {
      return this.checkDeleteAccess(USER_PERMISSIONS.technical.planningExceptions);
    },
  },
};
