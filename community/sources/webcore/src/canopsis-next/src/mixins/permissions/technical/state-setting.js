import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalStateSettingMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyStateSettingAccess() {
      return this.checkCreateAccess(USERS_PERMISSIONS.technical.shareToken); // TODO: change it
    },

    hasReadAnyStateSettingAccess() {
      return this.checkReadAccess(USERS_PERMISSIONS.technical.shareToken); // TODO: change it
    },

    hasUpdateAnyStateSettingAccess() {
      return this.checkUpdateAccess(USERS_PERMISSIONS.technical.shareToken); // TODO: change it
    },

    hasDeleteAnyStateSettingAccess() {
      return this.checkDeleteAccess(USERS_PERMISSIONS.technical.shareToken); // TODO: change it
    },
  },
};
