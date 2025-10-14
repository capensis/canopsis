import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Creates an instance of the dynamic info store module hooks.
 *
 * @returns {Object} Store module hooks for the 'dynamicInfo' namespace
 */
export const useDynamicInfoStoreModule = () => useStoreModuleHooks('dynamicInfo');

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
 *
 * @returns {Object} Hook return object
 * @property {DynamicInfoGetters} getters - All available getters
 * @property {DynamicInfoActions} actions - All available actions
 */
export const useDynamicInfo = () => {
  const { useGetters, useActions } = useDynamicInfoStoreModule();

  const getters = useGetters({
    dynamicInfosPending: 'pending',
    dynamicInfos: 'items',
    dynamicInfosMeta: 'meta',
  });

  const actions = useActions({
    fetchDynamicInfosList: 'fetchList',
    createDynamicInfo: 'create',
    updateDynamicInfo: 'update',
    removeDynamicInfo: 'remove',
    fetchDynamicInfosKeysWithoutStore: 'fetchInfosKeysWithoutStore',
  });

  return {
    ...getters,
    ...actions,
  };
};
