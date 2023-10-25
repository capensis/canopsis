import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsAlarmsListBookmark = {
  mixins: [authMixin],
  computed: {
    hasAccessToFilterByBookmark() {
      return this.checkAccess(USERS_PERMISSIONS.business.alarmsList.actions.filterByBookmark);
    },
  },
};
