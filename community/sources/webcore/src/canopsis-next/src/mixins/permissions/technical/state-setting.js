import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalStateSettingMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyStateSettingAccess() {
      return this.checkCreateAccess(USERS_PERMISSIONS.technical.stateSetting);
    },

    hasReadAnyStateSettingAccess() {
      return this.checkReadAccess(USERS_PERMISSIONS.technical.stateSetting);
    },

    hasUpdateAnyStateSettingAccess() {
      return this.checkUpdateAccess(USERS_PERMISSIONS.technical.stateSetting);
    },

    hasDeleteAnyStateSettingAccess() {
      return this.checkDeleteAccess(USERS_PERMISSIONS.technical.stateSetting);
    },
  },
};
