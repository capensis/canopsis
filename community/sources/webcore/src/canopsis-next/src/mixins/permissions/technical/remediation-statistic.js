import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalRemediationStatisticMixin = {
  mixins: [authMixin],
  computed: {
    hasReadRemediationStatisticAccess() {
      return this.checkAccess(USERS_PERMISSIONS.technical.remediationStatistic);
    },
  },
};
