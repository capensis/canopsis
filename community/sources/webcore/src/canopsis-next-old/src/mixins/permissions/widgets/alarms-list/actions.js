import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsAlarmsListActions = {
  mixins: [authMixin],
  computed: {
    hasAccessAlarmsListWidgetActions() {
      return this.checkAccess(USER_PERMISSIONS.business.alarmsList.actions.filterByBookmark);
    },
  },
};
