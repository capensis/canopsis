import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsAlarmsListCorrelation = {
  mixins: [authMixin],
  computed: {
    hasAccessToCorrelation() {
      return this.checkAccess(USER_PERMISSIONS.business.alarmsList.actions.correlation);
    },
  },
};
