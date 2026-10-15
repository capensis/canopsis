import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsContextPbehavior = {
  mixins: [authMixin],
  computed: {
    hasAccessToPbehavior() {
      return this.checkAccess(USER_PERMISSIONS.business.context.actions.pbehavior);
    },
  },
};
