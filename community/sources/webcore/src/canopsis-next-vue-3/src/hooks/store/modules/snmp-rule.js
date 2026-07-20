import { storeToRefs } from 'pinia'

import { useSnmpRuleStore } from '@/stores/modules/snmp-rule'

/**
 * SNMP rule store facade.
 *
 * @returns {Object}
 */
export const useSnmpRule = () => {
  const snmpRuleStore = useSnmpRuleStore()
  const { items, meta, pending, fetchingParams } = storeToRefs(snmpRuleStore)

  return {
    snmpRules: items,
    snmpRulesMeta: meta,
    snmpRulesPending: pending,
    snmpRulesFetchingParams: fetchingParams,

    fetchSnmpRulesList: snmpRuleStore.fetchList,
    fetchSnmpRulesListWithPreviousParams: snmpRuleStore.fetchListWithPreviousParams,
    fetchSnmpRulesListWithoutStore: snmpRuleStore.fetchListWithoutStore,
    createSnmpRule: snmpRuleStore.create,
    updateSnmpRule: snmpRuleStore.update,
    removeSnmpRule: snmpRuleStore.remove,
    bulkEnableSnmpRules: snmpRuleStore.bulkEnable,
    bulkDisableSnmpRules: snmpRuleStore.bulkDisable,
    bulkRemoveSnmpRules: snmpRuleStore.bulkRemove,
  }
}
