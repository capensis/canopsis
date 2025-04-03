import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Hook to access the SNMP MIB store module.
 *
 * @returns {Object} The store module hooks for the SNMP MIB.
 */
const useSnmpMibStoreModule = () => useStoreModuleHooks('snmpMib');

/**
 * Custom hook to use SNMP MIB actions.
 *
 * @returns {Object} An object containing SNMP MIB actions.
 * @property {Function} fetchSnmpMibList - Action to fetch the list of SNMP MIBs.
 */
export const useSnmpMib = () => {
  const { useActions } = useSnmpMibStoreModule();

  const actions = useActions({
    fetchSnmpMibList: 'fetchList',
  });

  return {
    ...actions,
  };
};
