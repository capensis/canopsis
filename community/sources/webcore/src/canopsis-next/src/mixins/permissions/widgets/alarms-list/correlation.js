import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsAlarmsListCorrelation = {
  mixins: [authMixin],
  computed: {
    hasAccessToCorrelation() {
      return this.checkAccess(USERS_PERMISSIONS.business.alarmsList.actions.correlation);
    },
  },
};
