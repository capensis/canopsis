import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalTagMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyTagAccess() {
      return this.checkCreateAccess(USER_PERMISSIONS.technical.tag);
    },

    hasReadAnyTagAccess() {
      return this.checkReadAccess(USER_PERMISSIONS.technical.tag);
    },

    hasUpdateAnyTagAccess() {
      return this.checkUpdateAccess(USER_PERMISSIONS.technical.tag);
    },

    hasDeleteAnyTagAccess() {
      return this.checkDeleteAccess(USER_PERMISSIONS.technical.tag);
    },
  },
};
