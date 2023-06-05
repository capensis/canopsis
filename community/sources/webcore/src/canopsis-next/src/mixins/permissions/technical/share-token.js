import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalShareTokenMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyShareTokenAccess() {
      return this.checkCreateAccess(USERS_PERMISSIONS.technical.shareToken);
    },

    hasReadAnyShareTokenAccess() {
      return this.checkReadAccess(USERS_PERMISSIONS.technical.shareToken);
    },

    hasUpdateAnyShareTokenAccess() {
      return this.checkUpdateAccess(USERS_PERMISSIONS.technical.shareToken);
    },

    hasDeleteAnyShareTokenAccess() {
      return this.checkDeleteAccess(USERS_PERMISSIONS.technical.shareToken);
    },
  },
};
