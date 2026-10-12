import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalIconMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyIconAccess() {
      return this.checkCreateAccess(USER_PERMISSIONS.technical.icon);
    },

    hasReadAnyIconAccess() {
      return this.checkReadAccess(USER_PERMISSIONS.technical.icon);
    },

    hasUpdateAnyIconAccess() {
      return this.checkUpdateAccess(USER_PERMISSIONS.technical.icon);
    },

    hasDeleteAnyIconAccess() {
      return this.checkDeleteAccess(USER_PERMISSIONS.technical.icon);
    },
  },
};
