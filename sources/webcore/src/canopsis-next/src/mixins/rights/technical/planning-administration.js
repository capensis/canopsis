import { USERS_RIGHTS } from '@/constants';

import authMixin from '@/mixins/auth';

export default {
  mixins: [authMixin],
  computed: {
    hasReadPlanningAdministrationAccess() {
      return this.checkReadAccess(USERS_RIGHTS.technical.planningAdministration);
    },

    hasUpdatePlanningAdministrationAccess() {
      return this.checkUpdateAccess(USERS_RIGHTS.technical.planningAdministration);
    },
  },
};
