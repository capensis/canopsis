import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalKpiCollectionSettingsMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyKpiCollectionSettingsAccess() {
      return this.checkCreateAccess(USERS_PERMISSIONS.technical.kpiCollectionSettings);
    },

    hasReadAnyKpiCollectionSettingsAccess() {
      return this.checkReadAccess(USERS_PERMISSIONS.technical.kpiCollectionSettings);
    },

    hasUpdateAnyKpiCollectionSettingsAccess() {
      return this.checkUpdateAccess(USERS_PERMISSIONS.technical.kpiCollectionSettings);
    },

    hasDeleteAnyKpiCollectionSettingsAccess() {
      return this.checkDeleteAccess(USERS_PERMISSIONS.technical.kpiCollectionSettings);
    },
  },
};
