import { isUndefined, omit } from 'lodash';

import uuid from '@/helpers/uuid';

/**
 * Convert a remediation instruction steps array to form array
 *
 * @param {RemediationInstructionStep[]} steps
 * @returns {Array}
 */
const remediationInstructionStepsToForm = steps => steps.map(step => ({
  ...step,
  operations: step.operations.map(operation => ({
    ...operation,
    saved: true,
    key: uuid(),
  })),
  saved: true,
  key: uuid(),
}));

/**
 * Convert a remediation instruction object to form object
 *
 * @typedef {Object} RemediationInstruction
 * @property {string} name
 * @property {boolean} enabled
 * @property {string} description
 * @property {Object} filter
 * @property {RemediationInstructionStep[]} steps
 * @param {Object} remediationInstruction
 * @returns {RemediationInstruction}
 */
export const remediationInstructionToForm = (remediationInstruction = {}) => ({
  name: remediationInstruction.name || '',
  enabled: !isUndefined(remediationInstruction.enabled) ? remediationInstruction.enabled : true,
  description: remediationInstruction.description || '',
  filter: remediationInstruction.filter || {},
  steps: remediationInstruction.steps
    ? remediationInstructionStepsToForm(remediationInstruction.steps)
    : [],
});


/**
 * Convert a remediation instruction steps form array to a API compatible array
 *
 * @param {RemediationInstructionStep[]} steps
 * @returns {Array}
 */
const formStepsToRemediationInstructionSteps = steps => steps.map(step => ({
  ...omit(step, ['key', 'saved']),
  operations: step.operations.map(operation => omit(operation, ['key', 'saved'])),
}));

/**
 * Convert a remediation instruction form object to a API compatible object
 *
 * @param {Object} form
 * @returns {Object}
 */
export const formToRemediationInstruction = form => ({
  ...form,
  steps: formStepsToRemediationInstructionSteps(form.steps),
});
