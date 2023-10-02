import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsUserStatisticsInterval = {
  mixins: [authMixin],
  computed: {
    hasAccessToInterval() {
      return this.checkAccess(USERS_PERMISSIONS.business.userStatistics.actions.interval);
    },
  },
};
