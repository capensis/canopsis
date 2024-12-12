import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsEventComment = {
  mixins: [authMixin],
  computed: {
    hasAccessToComment() {
      return this.checkAccess(USER_PERMISSIONS.business.context.actions.entityComment);
    },
  },
};
