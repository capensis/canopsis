import { USERS_RIGHTS } from '@/constants';

import authMixin from '@/mixins/auth';

export default {
  mixins: [authMixin],
  computed: {
    hasCreateAnyPbehaviorExceptionAccess() {
      return this.checkCreateAccess(USERS_RIGHTS.technical.exceptions);
    },

    hasReadAnyPbehaviorExceptionAccess() {
      return this.checkReadAccess(USERS_RIGHTS.technical.exceptions);
    },

    hasUpdateAnyPbehaviorExceptionAccess() {
      return this.checkUpdateAccess(USERS_RIGHTS.technical.exceptions);
    },

    hasDeleteAnyPbehaviorExceptionAccess() {
      return this.checkDeleteAccess(USERS_RIGHTS.technical.exceptions);
    },
  },
};
