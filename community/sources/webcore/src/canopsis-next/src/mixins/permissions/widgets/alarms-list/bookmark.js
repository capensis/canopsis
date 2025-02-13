import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsAlarmsListBookmark = {
  mixins: [authMixin],
  computed: {
    hasAccessToFilterByBookmark() {
      return this.checkAccess(USER_PERMISSIONS.business.alarmsList.actions.filterByBookmark);
    },
  },
};
