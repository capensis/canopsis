import { USERS_RIGHTS } from '@/constants';

import authMixin from '@/mixins/auth';

export default {
  mixins: [authMixin],
  computed: {
    hasCreateAnyPbehaviorExceptionAccess() {
      return this.checkCreateAccess(USERS_RIGHTS.technical.planningExceptions);
    },

    hasReadAnyPbehaviorExceptionAccess() {
      return this.checkReadAccess(USERS_RIGHTS.technical.planningExceptions);
    },

    hasUpdateAnyPbehaviorExceptionAccess() {
      return this.checkUpdateAccess(USERS_RIGHTS.technical.planningExceptions);
    },

    hasDeleteAnyPbehaviorExceptionAccess() {
      return this.checkDeleteAccess(USERS_RIGHTS.technical.planningExceptions);
    },
  },
};
