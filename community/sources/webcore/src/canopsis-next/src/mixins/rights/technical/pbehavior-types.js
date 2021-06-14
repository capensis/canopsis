import { USERS_RIGHTS } from '@/constants';

import authMixin from '@/mixins/auth';

export default {
  mixins: [authMixin],
  computed: {
    hasCreateAnyPbehaviorTypeAccess() {
      return this.checkCreateAccess(USERS_RIGHTS.technical.planningType);
    },

    hasReadAnyPbehaviorTypeAccess() {
      return this.checkReadAccess(USERS_RIGHTS.technical.planningType);
    },

    hasUpdateAnyPbehaviorTypeAccess() {
      return this.checkUpdateAccess(USERS_RIGHTS.technical.planningType);
    },

    hasDeleteAnyPbehaviorTypeAccess() {
      return this.checkDeleteAccess(USERS_RIGHTS.technical.planningType);
    },
  },
};
