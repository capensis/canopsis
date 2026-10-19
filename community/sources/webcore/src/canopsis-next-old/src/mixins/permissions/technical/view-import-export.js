import { USER_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsTechnicalViewImportExportMixin = {
  mixins: [authMixin],
  computed: {
    hasAccessToViewImportExport() {
      return this.checkAccess(USER_PERMISSIONS.technical.viewImportExport);
    },
  },
};
