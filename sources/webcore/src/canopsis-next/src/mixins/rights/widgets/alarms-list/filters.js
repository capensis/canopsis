import { USERS_RIGHTS } from '@/constants';

import authMixin from '@/mixins/auth';

export default {
  mixins: [authMixin],
  computed: {
    hasAccessToListFilters() {
      return this.checkAccess(USERS_RIGHTS.business.alarmsList.actions.listFilters);
    },

    hasAccessToEditFilter() {
      return this.checkAccess(USERS_RIGHTS.business.alarmsList.actions.editFilter);
    },

    hasAccessToAddFilter() {
      return this.checkAccess(USERS_RIGHTS.business.alarmsList.actions.addFilter);
    },

    hasAccessToUserFilter() {
      return this.checkAccess(USERS_RIGHTS.business.alarmsList.actions.userFilter);
    },
  },
};
