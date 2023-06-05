import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalProfileCorporatePatternMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyCorporatePatternAccess() {
      return this.checkCreateAccess(USERS_PERMISSIONS.technical.profile.corporatePattern);
    },

    hasReadAnyCorporatePatternAccess() {
      return this.checkReadAccess(USERS_PERMISSIONS.technical.profile.corporatePattern);
    },

    hasUpdateAnyCorporatePatternAccess() {
      return this.checkUpdateAccess(USERS_PERMISSIONS.technical.profile.corporatePattern);
    },

    hasDeleteAnyCorporatePatternAccess() {
      return this.checkDeleteAccess(USERS_PERMISSIONS.technical.profile.corporatePattern);
    },
  },
};
