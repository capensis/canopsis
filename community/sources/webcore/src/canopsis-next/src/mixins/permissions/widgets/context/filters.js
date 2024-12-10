import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsContextFilters = {
  mixins: [authMixin],
  computed: {
    hasAccessToListFilters() {
      return this.checkAccess(USER_PERMISSIONS.business.context.actions.listFilters);
    },

    hasAccessToEditFilter() {
      return this.checkAccess(USER_PERMISSIONS.business.context.actions.editFilter);
    },

    hasAccessToAddFilter() {
      return this.checkAccess(USER_PERMISSIONS.business.context.actions.addFilter);
    },

    hasAccessToUserFilter() {
      return this.checkAccess(USER_PERMISSIONS.business.context.actions.userFilter);
    },
  },
};
