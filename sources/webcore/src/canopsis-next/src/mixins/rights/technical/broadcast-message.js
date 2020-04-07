import { USERS_RIGHTS } from '@/constants';

import authMixin from '@/mixins/auth';

export default {
  mixins: [authMixin],
  computed: {
    hasCreateAnyBroadcastMessageAccess() {
      return this.checkCreateAccess(USERS_RIGHTS.technical.broadcastMessage);
    },

    hasReadAnyBroadcastMessageAccess() {
      return this.checkReadAccess(USERS_RIGHTS.technical.broadcastMessage);
    },

    hasUpdateAnyBroadcastMessageAccess() {
      return this.checkUpdateAccess(USERS_RIGHTS.technical.broadcastMessage);
    },

    hasDeleteAnyBroadcastMessageAccess() {
      return this.checkDeleteAccess(USERS_RIGHTS.technical.broadcastMessage);
    },
  },
};
