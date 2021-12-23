import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalKpiRatingSettingsMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyKpiRatingSettingsAccess() {
      return this.checkCreateAccess(USERS_PERMISSIONS.technical.kpiRatingSettings);
    },

    hasReadAnyKpiRatingSettingsAccess() {
      return this.checkReadAccess(USERS_PERMISSIONS.technical.kpiRatingSettings);
    },

    hasUpdateAnyKpiRatingSettingsAccess() {
      return this.checkUpdateAccess(USERS_PERMISSIONS.technical.kpiRatingSettings);
    },

    hasDeleteAnyKpiRatingSettingsAccess() {
      return this.checkDeleteAccess(USERS_PERMISSIONS.technical.kpiRatingSettings);
    },
  },
};
