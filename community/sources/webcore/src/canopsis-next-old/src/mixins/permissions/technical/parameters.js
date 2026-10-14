import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalParametersMixin = {
  mixins: [authMixin],
  computed: {
    hasReadParametersAccess() {
      return this.checkReadAccess(USER_PERMISSIONS.technical.parameters);
    },

    hasUpdateParametersAccess() {
      return this.checkUpdateAccess(USER_PERMISSIONS.technical.parameters);
    },
  },
};
