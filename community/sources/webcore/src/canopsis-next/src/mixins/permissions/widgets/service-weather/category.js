import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

export const permissionsWidgetsServiceWeatherCategory = {
  mixins: [authMixin],
  computed: {
    hasAccessToCategory() {
      return this.checkAccess(USERS_PERMISSIONS.business.serviceWeather.actions.category);
    },
  },
};
