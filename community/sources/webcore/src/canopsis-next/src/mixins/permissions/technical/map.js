import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalMapMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyMapAccess() {
      return this.checkCreateAccess(USERS_PERMISSIONS.technical.map);
    },

    hasReadAnyMapAccess() {
      return this.checkReadAccess(USERS_PERMISSIONS.technical.map);
    },

    hasUpdateAnyMapAccess() {
      return this.checkUpdateAccess(USERS_PERMISSIONS.technical.map);
    },

    hasDeleteAnyMapAccess() {
      return this.checkDeleteAccess(USERS_PERMISSIONS.technical.map);
    },
  },
};
