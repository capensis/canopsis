import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalWidgetTemplateMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyWidgetTemplateAccess() {
      return this.checkCreateAccess(USERS_PERMISSIONS.technical.widgetTemplate);
    },

    hasReadAnyWidgetTemplateAccess() {
      return this.checkReadAccess(USERS_PERMISSIONS.technical.widgetTemplate);
    },

    hasUpdateAnyWidgetTemplateAccess() {
      return this.checkUpdateAccess(USERS_PERMISSIONS.technical.widgetTemplate);
    },

    hasDeleteAnyWidgetTemplateAccess() {
      return this.checkDeleteAccess(USERS_PERMISSIONS.technical.widgetTemplate);
    },
  },
};
