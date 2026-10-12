import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Creates hooks for accessing the SNMP rule Vuex store module.
 * Provides access to getters and actions for managing SNMP rule operations.
 *
 * @returns {Object} An object containing store module utilities:
 * @property {import('vuex').Store} store - The Vuex store instance
 * @property {import('vuex').Module} module - The SNMP rule module instance
 * @property {Function} useGetters - Function to access module getters
 * @property {Function} useActions - Function to access module actions
 */
const useSnmpRuleStoreModule = () => useStoreModuleHooks('snmpRule');

/**
 * Custom hook for SNMP rule operations.
 * Provides convenient access to SNMP rule getters and actions.
 *
 * @returns {Object} An object containing SNMP rule getters and actions:
 * @property {Object} snmpRulesMeta - Getter for SNMP rule metadata
 * @property {boolean} snmpRulesPending - Getter for SNMP rule loading state
 * @property {Array} snmpRules - Getter for SNMP rule items
 * @property {Function} fetchSnmpRulesList - Action to fetch SNMP rules list
 * @property {Function} createSnmpRule - Action to create an SNMP rule
 * @property {Function} updateSnmpRule - Action to update an SNMP rule
 * @property {Function} removeSnmpRule - Action to remove an SNMP rule
 * @property {Function} bulkEnableSnmpRules - Action to bulk enable SNMP rules
 * @property {Function} bulkDisableSnmpRules - Action to bulk disable SNMP rules
 * @property {Function} bulkRemoveSnmpRules - Action to bulk remove SNMP rules
 * @property {Function} fetchSnmpRulesListWithPreviousParams - Action to fetch SNMP rules list with previous params
 * @property {Function} fetchSnmpRulesListWithoutStore - Action to fetch SNMP rules list without store
 *
 * @example
 * // Usage in a component
 * const { snmpRules, snmpRulesPending, fetchSnmpRulesList } = useSnmpRule();
 * await fetchSnmpRulesList({ params: { page: 1, limit: 10 } });
 */
export const useSnmpRule = () => {
  const { useGetters, useActions } = useSnmpRuleStoreModule();

  const getters = useGetters({
    snmpRulesMeta: 'meta',
    snmpRulesPending: 'pending',
    snmpRules: 'items',
  });

  const actions = useActions({
    fetchSnmpRulesList: 'fetchList',
    fetchSnmpRulesListWithPreviousParams: 'fetchListWithPreviousParams',
    fetchSnmpRulesListWithoutStore: 'fetchListWithoutStore',
    createSnmpRule: 'create',
    updateSnmpRule: 'update',
    removeSnmpRule: 'remove',
    bulkEnableSnmpRules: 'bulkEnable',
    bulkDisableSnmpRules: 'bulkDisable',
    bulkRemoveSnmpRules: 'bulkRemove',
  });

  return {
    ...getters,
    ...actions,
  };
};
