import { USERS_RIGHTS } from '@/constants';

import authMixin from '@/mixins/auth';

export default {
  mixins: [authMixin],
  computed: {
    hasCreateAnyPbehaviorTypeAccess() {
      return this.checkCreateAccess(USERS_RIGHTS.technical.type);
    },

    hasReadAnyPbehaviorTypeAccess() {
      return this.checkReadAccess(USERS_RIGHTS.technical.type);
    },

    hasUpdateAnyPbehaviorTypeAccess() {
      return this.checkUpdateAccess(USERS_RIGHTS.technical.type);
    },

    hasDeleteAnyPbehaviorTypeAccess() {
      return this.checkDeleteAccess(USERS_RIGHTS.technical.type);
    },
  },
};
