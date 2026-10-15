import { inject } from 'vue';

import { createPopupsApi } from '@/plugins/popups';
import { usePopupsStore } from '@/plugins/popups/stores';

/**
 * Get popups API.
 *
 * @return {Object}
 */
export const usePopups = () => {
  const providedPopups = inject('$popups', null);

  if (providedPopups) {
    return providedPopups;
  }

  return createPopupsApi(usePopupsStore());
};
