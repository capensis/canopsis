import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalIconMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyIconAccess() {
      return this.checkCreateAccess(USERS_PERMISSIONS.technical.icon);
    },

    hasReadAnyIconAccess() {
      return this.checkReadAccess(USERS_PERMISSIONS.technical.icon);
    },

    hasUpdateAnyIconAccess() {
      return this.checkUpdateAccess(USERS_PERMISSIONS.technical.icon);
    },

    hasDeleteAnyIconAccess() {
      return this.checkDeleteAccess(USERS_PERMISSIONS.technical.icon);
    },
  },
};
