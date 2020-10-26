import { USERS_RIGHTS } from '@/constants';

import authMixin from '@/mixins/auth';

export default {
  mixins: [authMixin],
  computed: {
    hasCreateAnyPbehaviorReasonAccess() {
      return this.checkCreateAccess(USERS_RIGHTS.technical.planningReason);
    },

    hasReadAnyPbehaviorReasonAccess() {
      return this.checkReadAccess(USERS_RIGHTS.technical.planningReason);
    },

    hasUpdateAnyPbehaviorReasonAccess() {
      return this.checkUpdateAccess(USERS_RIGHTS.technical.planningReason);
    },

    hasDeleteAnyPbehaviorReasonAccess() {
      return this.checkDeleteAccess(USERS_RIGHTS.technical.planningReason);
    },
  },
};
