import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalTagMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyTagAccess() {
      return this.checkCreateAccess(USERS_PERMISSIONS.technical.tag);
    },

    hasReadAnyTagAccess() {
      return this.checkReadAccess(USERS_PERMISSIONS.technical.tag);
    },

    hasUpdateAnyTagAccess() {
      return this.checkUpdateAccess(USERS_PERMISSIONS.technical.tag);
    },

    hasDeleteAnyTagAccess() {
      return this.checkDeleteAccess(USERS_PERMISSIONS.technical.tag);
    },
  },
};
