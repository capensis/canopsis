import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalPermissionMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyPermissionAccess() {
      return this.checkCreateAccess(USER_PERMISSIONS.technical.permission);
    },

    hasReadAnyPermissionAccess() {
      return this.checkReadAccess(USER_PERMISSIONS.technical.permission);
    },

    hasUpdateAnyPermissionAccess() {
      return this.checkUpdateAccess(USER_PERMISSIONS.technical.permission);
    },

    hasDeleteAnyPermissionAccess() {
      return this.checkDeleteAccess(USER_PERMISSIONS.technical.permission);
    },
  },
};
