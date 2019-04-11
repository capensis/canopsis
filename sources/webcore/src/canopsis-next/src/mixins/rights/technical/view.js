import { USERS_RIGHTS } from '@/constants';

import authMixin from '../../auth';

export default {
  mixins: [authMixin],
  computed: {
    hasCreateAnyViewAccess() {
      return this.checkCreateAccess(USERS_RIGHTS.technical.view);
    },

    hasReadAnyViewAccess() {
      return this.checkReadAccess(USERS_RIGHTS.technical.view);
    },

    hasUpdateAnyViewAccess() {
      return this.checkUpdateAccess(USERS_RIGHTS.technical.view);
    },

    hasDeleteAnyViewAccess() {
      return this.checkDeleteAccess(USERS_RIGHTS.technical.view);
    },
  },
};
