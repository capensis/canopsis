import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalBroadcastMessageMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyBroadcastMessageAccess() {
      return this.checkCreateAccess(USER_PERMISSIONS.technical.broadcastMessage);
    },

    hasReadAnyBroadcastMessageAccess() {
      return this.checkReadAccess(USER_PERMISSIONS.technical.broadcastMessage);
    },

    hasUpdateAnyBroadcastMessageAccess() {
      return this.checkUpdateAccess(USER_PERMISSIONS.technical.broadcastMessage);
    },

    hasDeleteAnyBroadcastMessageAccess() {
      return this.checkDeleteAccess(USER_PERMISSIONS.technical.broadcastMessage);
    },
  },
};
