import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalKpiCollectionSettingsMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyKpiCollectionSettingsAccess() {
      return this.checkCreateAccess(USER_PERMISSIONS.technical.kpiCollectionSettings);
    },

    hasReadAnyKpiCollectionSettingsAccess() {
      return this.checkReadAccess(USER_PERMISSIONS.technical.kpiCollectionSettings);
    },

    hasUpdateAnyKpiCollectionSettingsAccess() {
      return this.checkUpdateAccess(USER_PERMISSIONS.technical.kpiCollectionSettings);
    },

    hasDeleteAnyKpiCollectionSettingsAccess() {
      return this.checkDeleteAccess(USER_PERMISSIONS.technical.kpiCollectionSettings);
    },
  },
};
