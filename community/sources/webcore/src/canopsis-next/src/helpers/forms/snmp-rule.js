import { SNMP_STATE_TYPES, SNMP_TEMPLATE_STATE_STATES } from '@/constants';

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
 * @property {SnmpRuleModuleMib} [stateoid]
 * @property {string} [info]
 * @property {string} [minor]
 * @property {string} [major]
 * @property {string} [critical]
 */

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
 * @typedef {Object} SnmpRuleOidMib
 * @property {string} oid
 * @property {string} name
 */

/**
 * @typedef {Object} SnmpRuleOidForm
 * @property {string} moduleName
 * @property {SnmpRuleOidMib} mib
 */

/**
 * @typedef {Object} SnmpRuleForm
 * @property {SnmpRuleModuleMib} component
 * @property {SnmpRuleModuleMib} connector_name
 * @property {SnmpRuleModuleMib} output
 * @property {SnmpRuleModuleMib} resource
 * @property {SnmpRuleOidForm} oid
 * @property {SnmpRuleState} state
 */

/**
 * Convert snmp rule oid field to form
 *
 * @param {SnmpRuleOid} oid
 * @returns {SnmpRuleOidForm}
 */
export const snmpRuleOidToForm = (oid = {}) => ({
  moduleName: oid.moduleName ?? '',
  mib: {
    oid: oid.oid ?? '',
    name: oid.mibName ?? '',
  },
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
 * Convert snmp rule state to form
 *
 * @param {SnmpRuleState} state
 * @returns {SnmpRuleState}
 */
export const snmpRuleStateToForm = (state = {}) => {
  const type = state.type ?? SNMP_STATE_TYPES.simple;

  if (type === SNMP_STATE_TYPES.simple) {
    return {
      type,
      state: state.state,
    };
  }

  const additional = Object.values(SNMP_TEMPLATE_STATE_STATES).reduce((acc, value) => {
    acc[value] = state[value] ?? '';

    return acc;
  }, {});

  return {
    ...additional,

    type,
    state: state.state,
    stateoid: snmpRuleModuleMibToForm(state.stateoid),
  };
};

/**
 * Convert snmp rule to form
 *
 * @param {SnmpRule} snmpRule
 * @returns {SnmpRuleForm}
 */
export const snmpRuleToForm = (snmpRule = {}) => ({
  oid: snmpRuleOidToForm(snmpRule.oid),
  component: snmpRuleModuleMibToForm(snmpRule.component),
  connector_name: snmpRuleModuleMibToForm(snmpRule.connector_name),
  output: snmpRuleModuleMibToForm(snmpRule.output),
  resource: snmpRuleModuleMibToForm(snmpRule.resource),
  state: snmpRuleStateToForm(snmpRule.state),
});

/**
 * Convert oid form to snmp rule oid field
 *
 * @param {SnmpRuleOidForm} form
 * @returns {SnmpRuleOid}
 */
export const snmpRuleFormToOid = form => ({
  oid: form.mib.oid,
  mibName: form.mib.name,
  moduleName: form.moduleName,
});

/**
 * Convert form to snmp rule
 *
 * @param {SnmpRuleForm} form
 * @returns {SnmpRule}
 */
export const formToSnmpRule = form => ({
  ...form,

  oid: snmpRuleFormToOid(form.oid),
});
