import { USERS_RIGHTS } from '@/constants';

import authMixin from '@/mixins/auth';

export default {
  mixins: [authMixin],
  computed: {
    hasReadParametersAccess() {
      return this.checkReadAccess(USERS_RIGHTS.technical.parameters);
    },

    hasUpdateParametersAccess() {
      return this.checkUpdateAccess(USERS_RIGHTS.technical.parameters);
    },
  },
};
