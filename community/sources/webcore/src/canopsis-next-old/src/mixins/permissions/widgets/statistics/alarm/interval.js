import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsAlarmStatisticsInterval = {
  mixins: [authMixin],
  computed: {
    hasAccessToInterval() {
      return this.checkAccess(USER_PERMISSIONS.business.alarmStatistics.actions.interval);
    },
  },
};
