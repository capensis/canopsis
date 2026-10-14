import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalPlaylistMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyPlaylistAccess() {
      return this.checkCreateAccess(USER_PERMISSIONS.technical.playlist);
    },

    hasReadAnyPlaylistAccess() {
      return this.checkReadAccess(USER_PERMISSIONS.technical.playlist);
    },

    hasUpdateAnyPlaylistAccess() {
      return this.checkUpdateAccess(USER_PERMISSIONS.technical.playlist);
    },

    hasDeleteAnyPlaylistAccess() {
      return this.checkDeleteAccess(USER_PERMISSIONS.technical.playlist);
    },
  },
};
