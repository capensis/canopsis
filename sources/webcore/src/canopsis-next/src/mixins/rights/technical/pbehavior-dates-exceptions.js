import { USERS_RIGHTS } from '@/constants';

import authMixin from '@/mixins/auth';

export default {
  mixins: [authMixin],
  computed: {
    hasCreateAnyPbehaviorDateExceptionAccess() {
      return this.checkCreateAccess(USERS_RIGHTS.technical.datesExceptions);
    },

    hasReadAnyPbehaviorDateExceptionAccess() {
      return this.checkReadAccess(USERS_RIGHTS.technical.datesExceptions);
    },

    hasUpdateAnyPbehaviorDateExceptionAccess() {
      return this.checkUpdateAccess(USERS_RIGHTS.technical.datesExceptions);
    },

    hasDeleteAnyPbehaviorDateExceptionAccess() {
      return this.checkDeleteAccess(USERS_RIGHTS.technical.datesExceptions);
    },
  },
};
