import { isUndefined } from 'lodash';

/**
 * Convert a remediation instruction object to form object
 *
 * @param {Object} remediationInstruction
 * @returns {Object}
 */
export const remediationInstructionToForm = (remediationInstruction = {}) => ({
  name: remediationInstruction.name || '',
  enabled: !isUndefined(remediationInstruction.enabled) ? remediationInstruction.enabled : true,
  filter: {},
  steps: [],
});

/**
 * Convert a remediation instruction form object to a API compatible object
 *
 * @param {Object} form
 * @returns {Object}
 */
export const formToRemediationInstruction = form => form;
