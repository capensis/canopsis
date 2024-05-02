import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsAvailabilityExport = {
  mixins: [authMixin],
  computed: {
    hasAccessToExportAsCsv() {
      return this.checkAccess(USERS_PERMISSIONS.business.availability.actions.exportAsCsv);
    },
  },
};
