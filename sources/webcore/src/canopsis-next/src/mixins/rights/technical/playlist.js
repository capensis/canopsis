import { USERS_RIGHTS } from '@/constants';

import authMixin from '@/mixins/auth';

export default {
  mixins: [authMixin],
  computed: {
    hasCreateAnyPlaylistAccess() {
      return this.checkCreateAccess(USERS_RIGHTS.technical.playlist);
    },

    hasReadAnyPlaylistAccess() {
      return this.checkReadAccess(USERS_RIGHTS.technical.playlist);
    },

    hasUpdateAnyPlaylistAccess() {
      return this.checkUpdateAccess(USERS_RIGHTS.technical.playlist);
    },

    hasDeleteAnyPlaylistAccess() {
      return this.checkDeleteAccess(USERS_RIGHTS.technical.playlist);
    },
  },
};
