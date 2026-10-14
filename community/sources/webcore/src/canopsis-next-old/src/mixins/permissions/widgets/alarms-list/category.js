import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsAlarmsListCategory = {
  mixins: [authMixin],
  computed: {
    hasAccessToCategory() {
      return this.checkAccess(USER_PERMISSIONS.business.alarmsList.actions.category);
    },
  },
};
