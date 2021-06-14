import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalBroadcastMessageMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyBroadcastMessageAccess() {
      return this.checkCreateAccess(USERS_PERMISSIONS.technical.broadcastMessage);
    },

    hasReadAnyBroadcastMessageAccess() {
      return this.checkReadAccess(USERS_PERMISSIONS.technical.broadcastMessage);
    },

    hasUpdateAnyBroadcastMessageAccess() {
      return this.checkUpdateAccess(USERS_PERMISSIONS.technical.broadcastMessage);
    },

    hasDeleteAnyBroadcastMessageAccess() {
      return this.checkDeleteAccess(USERS_PERMISSIONS.technical.broadcastMessage);
    },
  },
};
