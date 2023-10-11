import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalProfileThemeMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyThemeAccess() {
      return this.checkCreateAccess(USERS_PERMISSIONS.technical.profile.theme);
    },

    hasReadAnyThemeAccess() {
      return this.checkReadAccess(USERS_PERMISSIONS.technical.profile.theme);
    },

    hasUpdateAnyThemeAccess() {
      return this.checkUpdateAccess(USERS_PERMISSIONS.technical.profile.theme);
    },

    hasDeleteAnyThemeAccess() {
      return this.checkDeleteAccess(USERS_PERMISSIONS.technical.profile.theme);
    },
  },
};
