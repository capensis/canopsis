import { ref } from 'vue';

import { API_ROUTES } from '@/config';
import request from '@/services/request';
import { createCRUDStore } from '@/stores/plugins/crud';

export const useRatingSettingsStore = createCRUDStore({
  id: 'ratingSettings',
  route: API_ROUTES.ratingSettings,
  withFetchingParams: true,
  withWithoutStore: true,
  setup() {
    const updatedAt = ref(null);

    const bulkUpdate = async ({ data } = {}) => {
      await request.put(API_ROUTES.bulkRatingSettings, data);

      updatedAt.value = Date.now();
    };

    return {
      updatedAt,
      bulkUpdate,
    };
  },
});
