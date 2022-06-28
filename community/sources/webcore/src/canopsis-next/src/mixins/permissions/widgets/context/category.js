import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsContextCategory = {
  mixins: [authMixin],
  computed: {
    hasAccessToCategory() {
      return this.checkAccess(USERS_PERMISSIONS.business.context.actions.category);
    },
  },
};
