import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsMapFilters = {
  mixins: [authMixin],
  computed: {
    hasAccessToListFilters() {
      return this.checkAccess(USER_PERMISSIONS.business.map.actions.listFilters);
    },

    hasAccessToEditFilter() {
      return this.checkAccess(USER_PERMISSIONS.business.map.actions.editFilter);
    },

    hasAccessToAddFilter() {
      return this.checkAccess(USER_PERMISSIONS.business.map.actions.addFilter);
    },

    hasAccessToUserFilter() {
      return this.checkAccess(USER_PERMISSIONS.business.map.actions.userFilter);
    },
  },
};
