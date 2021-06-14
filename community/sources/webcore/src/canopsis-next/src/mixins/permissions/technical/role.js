import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalRoleMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyRoleAccess() {
      return this.checkCreateAccess(USERS_PERMISSIONS.technical.role);
    },

    hasReadAnyRoleAccess() {
      return this.checkReadAccess(USERS_PERMISSIONS.technical.role);
    },

    hasUpdateAnyRoleAccess() {
      return this.checkUpdateAccess(USERS_PERMISSIONS.technical.role);
    },

    hasDeleteAnyRoleAccess() {
      return this.checkDeleteAccess(USERS_PERMISSIONS.technical.role);
    },
  },
};
