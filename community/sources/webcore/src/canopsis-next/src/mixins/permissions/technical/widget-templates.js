import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalWidgetTemplateMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyWidgetTemplateAccess() {
      return this.checkCreateAccess(USERS_PERMISSIONS.technical.view); // TODO: change permission
    },

    hasReadAnyWidgetTemplateAccess() {
      return this.checkReadAccess(USERS_PERMISSIONS.technical.view);
    },

    hasUpdateAnyWidgetTemplateAccess() {
      return this.checkUpdateAccess(USERS_PERMISSIONS.technical.view);
    },

    hasDeleteAnyWidgetTemplateAccess() {
      return this.checkDeleteAccess(USERS_PERMISSIONS.technical.view);
    },
  },
};
