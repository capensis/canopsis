import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalMapMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyMapAccess() {
      return this.checkCreateAccess(USER_PERMISSIONS.technical.map);
    },

    hasReadAnyMapAccess() {
      return this.checkReadAccess(USER_PERMISSIONS.technical.map);
    },

    hasUpdateAnyMapAccess() {
      return this.checkUpdateAccess(USER_PERMISSIONS.technical.map);
    },

    hasDeleteAnyMapAccess() {
      return this.checkDeleteAccess(USER_PERMISSIONS.technical.map);
    },
  },
};
