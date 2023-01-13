/**
 * @typedef {Object} SnmpRuleModuleMib
 * @property {string} formatter
 * @property {string} regex
 * @property {string} value
 */

/**
 * @typedef {Object} SnmpRuleOid
 * @property {string} mibName
 * @property {string} moduleName
 * @property {string} oid
 */

/**
 * @typedef {Object} SnmpRuleState
 * @property {string} type
 * @property {number} [state]
 */

import { SNMP_STATE_TYPES } from '@/constants';

/**
 * @typedef {Object} SnmpRule
 * @property {SnmpRuleModuleMib} component
 * @property {SnmpRuleModuleMib} connector_name
 * @property {SnmpRuleModuleMib} output
 * @property {SnmpRuleModuleMib} resource
 * @property {SnmpRuleOid} oid
 * @property {SnmpRuleState} state
 */

/**
 * Convert snmp rule oid field to form
 *
 * @param {SnmpRuleOid} oid
 * @returns {SnmpRuleOid}
 */
export const snmpRuleOidToForm = (oid = {}) => ({
  oid: oid.oid ?? '',
  mibName: oid.mibName ?? '',
  moduleName: oid.moduleName ?? '',
});

/**
 * Convert snmp rule module mib field to form
 *
 * @param {SnmpRuleModuleMib} moduleMib
 * @returns {SnmpRuleModuleMib}
 */
export const snmpRuleModuleMibToForm = (moduleMib = {}) => ({
  value: moduleMib.value ?? '',
  regex: moduleMib.regex ?? '',
  formatter: moduleMib.formatter ?? '',
});

/**
 * Convert snmp rule to form
 *
 * @param {SnmpRule} snmpRule
 * @returns {SnmpRule}
 */
export const snmpRuleToForm = (snmpRule = {}) => ({
  oid: snmpRuleOidToForm(snmpRule.oid),
  component: snmpRuleModuleMibToForm(snmpRule.component),
  connector_name: snmpRuleModuleMibToForm(snmpRule.connector_name),
  output: snmpRuleModuleMibToForm(snmpRule.output),
  resource: snmpRuleModuleMibToForm(snmpRule.resource),
  state: {
    type: snmpRule.state?.type ?? SNMP_STATE_TYPES.simple,
  },
});
