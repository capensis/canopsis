import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Retrieves the store module hooks for managing meta alarm rules.
 *
 * @returns {Object} An object containing store module.
 */
const useMetaAlarmRuleStoreModule = () => useStoreModuleHooks('metaAlarmRule');

/**
 * Provides actions to interact with meta alarm rules.
 *
 * @returns {Object} An object containing actions to create, update, remove, and fetch meta alarm rules list without
 * storing.
 */
export const useMetaAlarmRule = () => {
  const { useActions } = useMetaAlarmRuleStoreModule();

  const actions = useActions({
    createMetaAlarmRule: 'create',
    updateMetaAlarmRule: 'update',
    removeMetaAlarmRule: 'remove',
    fetchMetaAlarmRulesListWithoutStore: 'fetchListWithoutStore',
  });

  return {
    ...actions,
  };
};
