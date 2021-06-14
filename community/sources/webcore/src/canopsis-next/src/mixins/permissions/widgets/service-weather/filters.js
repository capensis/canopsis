import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsServiceWeatherFilters = {
  mixins: [authMixin],
  computed: {
    hasAccessToUserFilter() {
      return this.checkAccess(USERS_PERMISSIONS.business.serviceWeather.actions.userFilter);
    },
  },
};
