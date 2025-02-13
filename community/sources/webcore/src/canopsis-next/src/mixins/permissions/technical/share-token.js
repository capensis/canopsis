import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalShareTokenMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyShareTokenAccess() {
      return this.checkCreateAccess(USER_PERMISSIONS.technical.shareToken);
    },

    hasReadAnyShareTokenAccess() {
      return this.checkReadAccess(USER_PERMISSIONS.technical.shareToken);
    },

    hasUpdateAnyShareTokenAccess() {
      return this.checkUpdateAccess(USER_PERMISSIONS.technical.shareToken);
    },

    hasDeleteAnyShareTokenAccess() {
      return this.checkDeleteAccess(USER_PERMISSIONS.technical.shareToken);
    },
  },
};
