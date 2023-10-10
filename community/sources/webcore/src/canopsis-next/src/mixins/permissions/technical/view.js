import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalViewMixin = {
  mixins: [authMixin],
  computed: {
    hasAccessToPrivateView() {
      return this.checkAccess(USERS_PERMISSIONS.technical.privateView);
    },

    hasCreateAnyViewAccess() {
      return this.checkCreateAccess(USERS_PERMISSIONS.technical.view);
    },

    hasReadAnyViewAccess() {
      return this.checkReadAccess(USERS_PERMISSIONS.technical.view);
    },

    hasUpdateAnyViewAccess() {
      return this.checkUpdateAccess(USERS_PERMISSIONS.technical.view);
    },

    hasDeleteAnyViewAccess() {
      return this.checkDeleteAccess(USERS_PERMISSIONS.technical.view);
    },
  },
};
