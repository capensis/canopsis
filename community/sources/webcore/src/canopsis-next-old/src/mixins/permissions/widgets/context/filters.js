import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsContextFilters = {
  mixins: [authMixin],
  computed: {
    hasAccessToFilter() {
      return this.checkAccess(USER_PERMISSIONS.business.context.actions.filter);
    },

    hasAccessToUserFilter() {
      return this.checkAccess(USER_PERMISSIONS.business.context.actions.userFilter);
    },
  },
};
