import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalPermissionMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyPermissionAccess() {
      return this.checkCreateAccess(USERS_PERMISSIONS.technical.permission);
    },

    hasReadAnyPermissionAccess() {
      return this.checkReadAccess(USERS_PERMISSIONS.technical.permission);
    },

    hasUpdateAnyPermissionAccess() {
      return this.checkUpdateAccess(USERS_PERMISSIONS.technical.permission);
    },

    hasDeleteAnyPermissionAccess() {
      return this.checkDeleteAccess(USERS_PERMISSIONS.technical.permission);
    },
  },
};
