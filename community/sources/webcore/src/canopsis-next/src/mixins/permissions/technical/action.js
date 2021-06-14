import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalActionMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyActionAccess() {
      return this.checkCreateAccess(USERS_PERMISSIONS.technical.action);
    },

    hasReadAnyActionAccess() {
      return this.checkReadAccess(USERS_PERMISSIONS.technical.action);
    },

    hasUpdateAnyActionAccess() {
      return this.checkUpdateAccess(USERS_PERMISSIONS.technical.action);
    },

    hasDeleteAnyActionAccess() {
      return this.checkDeleteAccess(USERS_PERMISSIONS.technical.action);
    },
  },
};
