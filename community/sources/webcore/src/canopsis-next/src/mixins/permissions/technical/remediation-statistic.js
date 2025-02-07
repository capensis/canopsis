import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalRemediationStatisticMixin = {
  mixins: [authMixin],
  computed: {
    hasReadRemediationStatisticAccess() {
      return this.checkAccess(USER_PERMISSIONS.technical.remediationStatistic);
    },
  },
};
