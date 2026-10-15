import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsNumbersSampling = {
  mixins: [authMixin],
  computed: {
    hasAccessToSampling() {
      return this.checkAccess(USER_PERMISSIONS.business.pieChart.actions.sampling);
    },
  },
};
