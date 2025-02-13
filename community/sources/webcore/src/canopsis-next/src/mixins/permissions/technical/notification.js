import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalNotificationMixin = {
  mixins: [authMixin],
  computed: {
    hasAccessToNotificationsSettings() {
      return this.checkAccess(USER_PERMISSIONS.technical.notification.common);
    },
  },
};
