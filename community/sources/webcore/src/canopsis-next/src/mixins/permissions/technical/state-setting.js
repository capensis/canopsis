import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalStateSettingMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyStateSettingAccess() {
      return this.checkCreateAccess(USER_PERMISSIONS.technical.stateSetting);
    },

    hasReadAnyStateSettingAccess() {
      return this.checkReadAccess(USER_PERMISSIONS.technical.stateSetting);
    },

    hasUpdateAnyStateSettingAccess() {
      return this.checkUpdateAccess(USER_PERMISSIONS.technical.stateSetting);
    },

    hasDeleteAnyStateSettingAccess() {
      return this.checkDeleteAccess(USER_PERMISSIONS.technical.stateSetting);
    },
  },
};
