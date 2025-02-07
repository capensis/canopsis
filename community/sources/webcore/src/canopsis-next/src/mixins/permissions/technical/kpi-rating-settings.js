import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalKpiRatingSettingsMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyKpiRatingSettingsAccess() {
      return this.checkCreateAccess(USER_PERMISSIONS.technical.kpiRatingSettings);
    },

    hasReadAnyKpiRatingSettingsAccess() {
      return this.checkReadAccess(USER_PERMISSIONS.technical.kpiRatingSettings);
    },

    hasUpdateAnyKpiRatingSettingsAccess() {
      return this.checkUpdateAccess(USER_PERMISSIONS.technical.kpiRatingSettings);
    },

    hasDeleteAnyKpiRatingSettingsAccess() {
      return this.checkDeleteAccess(USER_PERMISSIONS.technical.kpiRatingSettings);
    },
  },
};
