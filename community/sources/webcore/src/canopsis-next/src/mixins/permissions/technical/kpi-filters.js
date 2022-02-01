import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalKpiFiltersMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyKpiFiltersAccess() {
      return this.checkCreateAccess(USERS_PERMISSIONS.technical.kpiFilters);
    },

    hasReadAnyKpiFiltersAccess() {
      return this.checkReadAccess(USERS_PERMISSIONS.technical.kpiFilters);
    },

    hasUpdateAnyKpiFiltersAccess() {
      return this.checkUpdateAccess(USERS_PERMISSIONS.technical.kpiFilters);
    },

    hasDeleteAnyKpiFiltersAccess() {
      return this.checkDeleteAccess(USERS_PERMISSIONS.technical.kpiFilters);
    },
  },
};
