import { USERS_RIGHTS } from '@/constants';

import authMixin from '@/mixins/auth';

export default {
  mixins: [authMixin],
  computed: {
    hasCreateAnyActionAccess() {
      return this.checkCreateAccess(USERS_RIGHTS.technical.action);
    },

    hasReadAnyActionAccess() {
      return this.checkReadAccess(USERS_RIGHTS.technical.action);
    },

    hasUpdateAnyActionAccess() {
      return this.checkUpdateAccess(USERS_RIGHTS.technical.action);
    },

    hasDeleteAnyActionAccess() {
      return this.checkDeleteAccess(USERS_RIGHTS.technical.action);
    },
  },
};
