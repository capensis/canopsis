import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalUserMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyUserAccess() {
      return this.checkCreateAccess(USERS_PERMISSIONS.technical.user);
    },

    hasReadAnyUserAccess() {
      return this.checkReadAccess(USERS_PERMISSIONS.technical.user);
    },

    hasUpdateAnyUserAccess() {
      return this.checkUpdateAccess(USERS_PERMISSIONS.technical.user);
    },

    hasDeleteAnyUserAccess() {
      return this.checkDeleteAccess(USERS_PERMISSIONS.technical.user);
    },
  },
};
