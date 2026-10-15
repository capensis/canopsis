import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalProfileThemeMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyThemeAccess() {
      return this.checkCreateAccess(USER_PERMISSIONS.technical.profile.theme);
    },

    hasReadAnyThemeAccess() {
      return this.checkReadAccess(USER_PERMISSIONS.technical.profile.theme);
    },

    hasUpdateAnyThemeAccess() {
      return this.checkUpdateAccess(USER_PERMISSIONS.technical.profile.theme);
    },

    hasDeleteAnyThemeAccess() {
      return this.checkDeleteAccess(USER_PERMISSIONS.technical.profile.theme);
    },
  },
};
