import { inject } from 'vue';

import { usePopupsStore } from '@/plugins/popups/stores';

import { createPopupsApi } from '@/plugins/popups';

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
