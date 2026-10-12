import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalWidgetTemplateMixin = {
  mixins: [authMixin],
  computed: {
    hasCreateAnyWidgetTemplateAccess() {
      return this.checkCreateAccess(USER_PERMISSIONS.technical.widgetTemplate);
    },

    hasReadAnyWidgetTemplateAccess() {
      return this.checkReadAccess(USER_PERMISSIONS.technical.widgetTemplate);
    },

    hasUpdateAnyWidgetTemplateAccess() {
      return this.checkUpdateAccess(USER_PERMISSIONS.technical.widgetTemplate);
    },

    hasDeleteAnyWidgetTemplateAccess() {
      return this.checkDeleteAccess(USER_PERMISSIONS.technical.widgetTemplate);
    },
  },
};
