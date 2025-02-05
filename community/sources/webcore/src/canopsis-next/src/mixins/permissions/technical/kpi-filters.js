import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalKpiFiltersMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyKpiFiltersAccess() {
      return this.checkCreateAccess(USER_PERMISSIONS.technical.kpiFilters);
    },

    hasReadAnyKpiFiltersAccess() {
      return this.checkReadAccess(USER_PERMISSIONS.technical.kpiFilters);
    },

    hasUpdateAnyKpiFiltersAccess() {
      return this.checkUpdateAccess(USER_PERMISSIONS.technical.kpiFilters);
    },

    hasDeleteAnyKpiFiltersAccess() {
      return this.checkDeleteAccess(USER_PERMISSIONS.technical.kpiFilters);
    },
  },
};
