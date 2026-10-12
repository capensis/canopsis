import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsAvailabilityFilters = {
  mixins: [authMixin],
  computed: {
    hasAccessToFilter() {
      return this.checkAccess(USER_PERMISSIONS.business.availability.actions.filter);
    },

    hasAccessToUserFilter() {
      return this.checkAccess(USER_PERMISSIONS.business.availability.actions.userFilter);
    },
  },
};
