import { storeToRefs } from 'pinia';

import { useDynamicInfoStore } from '@/stores/modules/dynamic-info';

/**
 * Hook for accessing and managing dynamic infos.
 * Provides CRUD operations and utility actions for dynamic information.
 *
 * @typedef {Object} DynamicInfoGetters
 * @property {Ref<boolean>} dynamicInfosPending - Loading state for dynamic infos
 * @property {Ref<Array|Object>} dynamicInfos - Dynamic infos items
 * @property {Ref<Object>} dynamicInfosMeta - Metadata for pagination
 *
 * @typedef {Object} DynamicInfoActions
 * @property {Function} fetchDynamicInfosList - Fetches dynamic infos list
 * @property {Function} createDynamicInfo - Creates a new dynamic info
 * @property {Function} updateDynamicInfo - Updates a dynamic info
 * @property {Function} removeDynamicInfo - Removes a dynamic info
 * @property {Function} fetchDynamicInfosKeysWithoutStore - Fetches dictionary keys without storing in Vuex
 * @property {Function} bulkEnableDynamicInfos - Bulk enables dynamic infos
 * @property {Function} bulkDisableDynamicInfos - Bulk disables dynamic infos
 * @property {Function} bulkRemoveDynamicInfos - Bulk removes dynamic infos
 *
 * @returns {DynamicInfoGetters & DynamicInfoActions}
 */
export const useDynamicInfo = () => {
  const dynamicInfoStore = useDynamicInfoStore();
  const {
    items,
    meta,
    pending,
    fetchingParams,
  } = storeToRefs(dynamicInfoStore);

  return {
    dynamicInfos: items,
    dynamicInfosMeta: meta,
    dynamicInfosPending: pending,
    dynamicInfosFetchingParams: fetchingParams,

    fetchDynamicInfosList: dynamicInfoStore.fetchList,
    fetchDynamicInfosListWithPreviousParams: dynamicInfoStore.fetchListWithPreviousParams,
    fetchDynamicInfosListWithoutStore: dynamicInfoStore.fetchListWithoutStore,
    createDynamicInfo: dynamicInfoStore.create,
    updateDynamicInfo: dynamicInfoStore.update,
    removeDynamicInfo: dynamicInfoStore.remove,
    fetchDynamicInfosKeysWithoutStore: dynamicInfoStore.fetchInfosKeysWithoutStore,
    bulkEnableDynamicInfos: dynamicInfoStore.bulkEnable,
    bulkDisableDynamicInfos: dynamicInfoStore.bulkDisable,
    bulkRemoveDynamicInfos: dynamicInfoStore.bulkRemove,
  };
};
