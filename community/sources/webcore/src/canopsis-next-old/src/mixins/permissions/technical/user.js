import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalUserMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyUserAccess() {
      return this.checkCreateAccess(USER_PERMISSIONS.technical.user);
    },

    hasReadAnyUserAccess() {
      return this.checkReadAccess(USER_PERMISSIONS.technical.user);
    },

    hasUpdateAnyUserAccess() {
      return this.checkUpdateAccess(USER_PERMISSIONS.technical.user);
    },

    hasDeleteAnyUserAccess() {
      return this.checkDeleteAccess(USER_PERMISSIONS.technical.user);
    },
  },
};
