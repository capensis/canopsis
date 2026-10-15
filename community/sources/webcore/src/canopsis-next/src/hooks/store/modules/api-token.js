import { storeToRefs } from 'pinia';

import { useApiTokenStore } from '@/stores/modules/api-token';

/**
 * Public API token store facade.
 *
 * @returns {Object}
 */
export const useApiToken = () => {
  const apiTokenStore = useApiTokenStore();
  const { items, meta, pending, fetchingParams } = storeToRefs(apiTokenStore);

  return {
    apiTokens: items,
    apiTokensMeta: meta,
    apiTokensPending: pending,
    apiTokensFetchingParams: fetchingParams,

    fetchApiTokensList: apiTokenStore.fetchList,
    fetchApiTokensListWithPreviousParams: apiTokenStore.fetchListWithPreviousParams,
    createApiToken: apiTokenStore.create,
    removeApiToken: apiTokenStore.remove,
  };
};
