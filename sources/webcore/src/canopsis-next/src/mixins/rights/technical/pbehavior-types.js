import { USERS_RIGHTS } from '@/constants';

import authMixin from '@/mixins/auth';

export default {
  mixins: [authMixin],
  computed: {
    hasCreateAnyTypeAccess() {
      return this.checkCreateAccess(USERS_RIGHTS.technical.type);
    },

    hasReadAnyTypeAccess() {
      return this.checkReadAccess(USERS_RIGHTS.technical.type);
    },

    hasUpdateAnyTypeAccess() {
      return this.checkUpdateAccess(USERS_RIGHTS.technical.type);
    },

    hasDeleteAnyTypeAccess() {
      return this.checkDeleteAccess(USERS_RIGHTS.technical.type);
    },
  },
};
