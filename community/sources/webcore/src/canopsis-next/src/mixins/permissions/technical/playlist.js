import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalPlaylistMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyPlaylistAccess() {
      return this.checkCreateAccess(USERS_PERMISSIONS.technical.playlist);
    },

    hasReadAnyPlaylistAccess() {
      return this.checkReadAccess(USERS_PERMISSIONS.technical.playlist);
    },

    hasUpdateAnyPlaylistAccess() {
      return this.checkUpdateAccess(USERS_PERMISSIONS.technical.playlist);
    },

    hasDeleteAnyPlaylistAccess() {
      return this.checkDeleteAccess(USERS_PERMISSIONS.technical.playlist);
    },
  },
};
