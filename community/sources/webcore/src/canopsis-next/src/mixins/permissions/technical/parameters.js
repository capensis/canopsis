import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalParametersMixin = {
  mixins: [authMixin],
  computed: {
    hasReadParametersAccess() {
      return this.checkReadAccess(USERS_PERMISSIONS.technical.parameters);
    },

    hasUpdateParametersAccess() {
      return this.checkUpdateAccess(USERS_PERMISSIONS.technical.parameters);
    },
  },
};
