import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsAvailabilityFilters = {
  mixins: [authMixin],
  computed: {
    hasAccessToListFilters() {
      return this.checkAccess(USERS_PERMISSIONS.business.availability.actions.listFilters);
    },

    hasAccessToEditFilter() {
      return this.checkAccess(USERS_PERMISSIONS.business.availability.actions.editFilter);
    },

    hasAccessToAddFilter() {
      return this.checkAccess(USERS_PERMISSIONS.business.availability.actions.addFilter);
    },

    hasAccessToUserFilter() {
      return this.checkAccess(USERS_PERMISSIONS.business.availability.actions.userFilter);
    },
  },
};
