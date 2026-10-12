import { SNMP_STATE_TYPES, SNMP_TEMPLATE_STATE_STATES } from '@/constants';

import { uid } from '@/helpers/uid';
import { removeKeyFromEntities } from '@/helpers/array';

/**
 * @typedef {Object} SnmpRuleModuleMib
 * @property {string} formatter
 * @property {string} regex
 * @property {string} value
 */

/**
 * @typedef {Object} SnmpRuleModuleExtra
 * @property {string} name
 * @property {SnmpRuleModuleMib} value
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
 * @property {boolean} enabled
 * @property {SnmpRuleModuleMib} component
 * @property {SnmpRuleModuleMib} connector_name
 * @property {SnmpRuleModuleMib} output
 * @property {SnmpRuleModuleMib} resource
 * @property {SnmpRuleModuleMib[]} tags
 * @property {SnmpRuleModuleExtra[]} extra
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
 * @typedef {SnmpRuleModuleExtra} SnmpRuleModuleExtraForm
 * @property {string} key
 */

/**
 * @typedef {SnmpRuleModuleMib} SnmpRuleModuleTagForm
 * @property {string} key
 */

/**
 * @typedef {Object} SnmpRuleForm
 * @property {boolean} enabled
 * @property {SnmpRuleModuleMib} component
 * @property {SnmpRuleModuleMib} connector_name
 * @property {SnmpRuleModuleMib} output
 * @property {SnmpRuleModuleMib} resource
 * @property {SnmpRuleModuleTagForm[]} tags
 * @property {SnmpRuleModuleExtraForm[]} extra
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
 * Convert snmp rule tag item to form format
 *
 * @param {SnmpRuleModuleMib} [tag={}] - The list of SNMP rule tags.
 * @returns {SnmpRuleModuleTagForm}
 */
export const snmpRuleTagToForm = (tag = {}) => ({
  ...snmpRuleModuleMibToForm(tag),

  key: uid(),
});

/**
 * Convert snmp rule tags to form format
 *
 * @param {SnmpRuleModuleMib[]} [tags=[]] - The list of SNMP rule tags.
 * @returns {SnmpRuleModuleTagForm[]}
 */
export const snmpRuleTagsToForm = (tags = []) => (tags ?? []).map(snmpRuleTagToForm);

/**
 * Convert snmp rule extra data item to form format
 *
 * @param {SnmpRuleModuleExtra} [extraItem={}] - The list of SNMP rule extra data.
 * @returns {SnmpRuleModuleExtraForm}
 */
export const snmpRuleExtraItemToForm = (extraItem = {}) => ({
  key: uid(),
  name: extraItem.name ?? '',
  value: snmpRuleModuleMibToForm(extraItem.value),
});

/**
 * Convert snmp rule extra data to form format
 *
 * @param {SnmpRuleModuleExtra[]} [extra=[]] - The list of SNMP rule extra data.
 * @returns {SnmpRuleModuleExtraForm[]}
 */
export const snmpRuleExtraToForm = (extra = []) => (extra ?? []).map(snmpRuleExtraItemToForm);

/**
 * Convert snmp rule to form
 *
 * @param {SnmpRule} snmpRule
 * @returns {SnmpRuleForm}
 */
export const snmpRuleToForm = (snmpRule = {}) => ({
  enabled: snmpRule.enabled ?? true,
  oid: snmpRuleOidToForm(snmpRule.oid),
  component: snmpRuleModuleMibToForm(snmpRule.component),
  connector_name: snmpRuleModuleMibToForm(snmpRule.connector_name),
  output: snmpRuleModuleMibToForm(snmpRule.output),
  resource: snmpRuleModuleMibToForm(snmpRule.resource),
  state: snmpRuleStateToForm(snmpRule.state),
  tags: snmpRuleTagsToForm(snmpRule.tags),
  extra: snmpRuleExtraToForm(snmpRule.extra),
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
  tags: removeKeyFromEntities(form.tags),
  extra: removeKeyFromEntities(form.extra),
});
