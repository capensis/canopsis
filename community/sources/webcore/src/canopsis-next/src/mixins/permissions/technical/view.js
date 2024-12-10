import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalViewMixin = {
  mixins: [authMixin],
  computed: {
    hasAccessToPrivateView() {
      return this.checkAccess(USER_PERMISSIONS.technical.privateView);
    },

    hasCreateAnyViewAccess() {
      return this.checkCreateAccess(USER_PERMISSIONS.technical.view);
    },

    hasReadAnyViewAccess() {
      return this.checkReadAccess(USER_PERMISSIONS.technical.view);
    },

    hasUpdateAnyViewAccess() {
      return this.checkUpdateAccess(USER_PERMISSIONS.technical.view);
    },

    hasDeleteAnyViewAccess() {
      return this.checkDeleteAccess(USER_PERMISSIONS.technical.view);
    },
  },
};
