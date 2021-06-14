import { USERS_RIGHTS } from '@/constants';

import authMixin from '@/mixins/auth';

export default {
  mixins: [authMixin],
  computed: {
    hasCreateAnyUserAccess() {
      return this.checkCreateAccess(USERS_RIGHTS.technical.user);
    },

    hasReadAnyUserAccess() {
      return this.checkReadAccess(USERS_RIGHTS.technical.user);
    },

    hasUpdateAnyUserAccess() {
      return this.checkUpdateAccess(USERS_RIGHTS.technical.user);
    },

    hasDeleteAnyUserAccess() {
      return this.checkDeleteAccess(USERS_RIGHTS.technical.user);
    },
  },
};
