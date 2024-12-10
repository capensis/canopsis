import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsEventComment = {
  mixins: [authMixin],
  computed: {
    hasAccessToCommentsList() {
      return this.checkAccess(USER_PERMISSIONS.business.context.actions.entityCommentsList);
    },

    hasAccessToCreateComment() {
      return this.checkAccess(USER_PERMISSIONS.business.context.actions.createEntityComment);
    },

    hasAccessToEditComment() {
      return this.checkAccess(USER_PERMISSIONS.business.context.actions.editEntityComment);
    },
  },
};
