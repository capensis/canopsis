import { isUndefined, omit } from 'lodash';

import uuid from '@/helpers/uuid';
import { addKeyInEntity, removeKeyFromEntity } from '@/helpers/entities';

/**
 * Convert a remediation instruction steps array to form array
 *
 * @param {Array} steps
 * @return {Array}
 */
const remediationInstructionStepsToForm = steps => steps.map(step => ({
  ...step,
  operations: addKeyInEntity(step.operations),
  saved: true,
  key: uuid(),
}));

/**
 * Convert a remediation instruction object to form object
 *
 * @param {Object} remediationInstruction
 * @returns {Object}
 */
export const remediationInstructionToForm = (remediationInstruction = {}) => ({
  name: remediationInstruction.name || '',
  enabled: !isUndefined(remediationInstruction.enabled) ? remediationInstruction.enabled : true,
  description: remediationInstruction.description || '',
  filter: remediationInstruction.filter || {},
  steps: remediationInstruction.steps ? remediationInstructionStepsToForm(remediationInstruction.steps) : [],
});


/**
 * Convert a remediation instruction steps form array to a API compatible array
 *
 * @param {Array} steps
 * @return {Array}
 */
const formStepsToRemediationInstructionSteps = steps => steps.map(step => ({
  ...omit(step, ['key', 'saved']),
  operations: removeKeyFromEntity(step.operations),
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
