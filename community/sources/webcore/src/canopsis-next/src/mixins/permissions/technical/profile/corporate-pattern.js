import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalProfileCorporatePatternMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyCorporatePatternAccess() {
      return this.checkCreateAccess(USER_PERMISSIONS.technical.profile.corporatePattern);
    },

    hasReadAnyCorporatePatternAccess() {
      return this.checkReadAccess(USER_PERMISSIONS.technical.profile.corporatePattern);
    },

    hasUpdateAnyCorporatePatternAccess() {
      return this.checkUpdateAccess(USER_PERMISSIONS.technical.profile.corporatePattern);
    },

    hasDeleteAnyCorporatePatternAccess() {
      return this.checkDeleteAccess(USER_PERMISSIONS.technical.profile.corporatePattern);
    },
  },
};
