import { useStoreModuleHooks } from '@/hooks/store';

const useSnmpRuleStoreModule = () => useStoreModuleHooks('snmpRule');

/**
 * Hook for accessing SNMP rule store module actions
 *
 * @returns {Object} An object containing SNMP rule actions
 */
export const useSnmpRule = () => {
  const { useActions } = useSnmpRuleStoreModule();

  const actions = useActions({
    fetchSnmpRulesList: 'fetchList',
    fetchSnmpRulesListWithPreviousParams: 'fetchListWithPreviousParams',
    fetchSnmpRulesListWithoutStore: 'fetchListWithoutStore',
    createSnmpRule: 'create',
    updateSnmpRule: 'update',
    removeSnmpRule: 'remove',
    bulkEnableSnmpRules: 'bulkEnable',
    bulkDisableSnmpRules: 'bulkDisable',
    bulkRemoveSnmpRules: 'bulkRemove', // TODO: check if available
  });

  return {
    ...actions,
  };
};
