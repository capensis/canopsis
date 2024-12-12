import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsMapFilters = {
  mixins: [authMixin],
  computed: {
    hasAccessToFilter() {
      return this.checkAccess(USER_PERMISSIONS.business.map.actions.filter);
    },

    hasAccessToUserFilter() {
      return this.checkAccess(USER_PERMISSIONS.business.map.actions.userFilter);
    },
  },
};
