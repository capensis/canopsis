import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalRoleMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyRoleAccess() {
      return this.checkCreateAccess(USER_PERMISSIONS.technical.role);
    },

    hasReadAnyRoleAccess() {
      return this.checkReadAccess(USER_PERMISSIONS.technical.role);
    },

    hasUpdateAnyRoleAccess() {
      return this.checkUpdateAccess(USER_PERMISSIONS.technical.role);
    },

    hasDeleteAnyRoleAccess() {
      return this.checkDeleteAccess(USER_PERMISSIONS.technical.role);
    },
  },
};
