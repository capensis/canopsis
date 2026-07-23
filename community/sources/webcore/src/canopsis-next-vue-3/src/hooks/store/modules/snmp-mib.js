import { useSnmpMibStore } from '@/stores/modules/snmp-mib';

/**
 * SNMP MIB store facade.
 *
 * @returns {Object}
 */
export const useSnmpMib = () => {
  const snmpMibStore = useSnmpMibStore();

  return {
    fetchSnmpMibsList: snmpMibStore.fetchList,
    uploadSnmpMib: snmpMibStore.upload,
  };
};
