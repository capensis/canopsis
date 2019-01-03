import { USERS_RIGHTS } from '@/constants';

import authMixin from '../../auth';

export default {
  mixins: [authMixin],
  computed: {
    hasCreateAnyRoleAccess() {
      return this.checkCreateAccess(USERS_RIGHTS.technical.role);
    },

    hasReadAnyRoleAccess() {
      return this.checkReadAccess(USERS_RIGHTS.technical.role);
    },

    hasUpdateAnyRoleAccess() {
      return this.checkUpdateAccess(USERS_RIGHTS.technical.role);
    },

    hasDeleteAnyRoleAccess() {
      return this.checkDeleteAccess(USERS_RIGHTS.technical.role);
    },
  },
};
