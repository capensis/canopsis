import { USERS_RIGHTS } from '@/constants';

import authMixin from '@/mixins/auth';

export default {
  mixins: [authMixin],
  computed: {
    hasAccessToCorrelation() {
      return this.checkAccess(USERS_RIGHTS.business.alarmsList.actions.correlation);
    },
  },
};
