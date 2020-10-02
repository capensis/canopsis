import { isUndefined, omit } from 'lodash';

import { TIME_UNITS } from '@/constants';

import { getSecondsByUnit, getUnitValueFromOtherUnit } from '@/helpers/time';
import uuid from '@/helpers/uuid';

/**
 * Convert a remediation instruction step operation array to form array
 *
 * @param {RemediationInstructionStepOperation[]} operations
 * @returns {Array}
 */
const remediationInstructionStepOperationsToForm = operations => operations.map(operation => ({
  ...operation,
  time_to_complete: {
    interval: getUnitValueFromOtherUnit(operation.time_to_complete, TIME_UNITS.second, operation.time_to_complete_unit),
    unit: operation.time_to_complete_unit || TIME_UNITS.second,
  },
  saved: true,
  key: uuid(),
}));

/**
 * Convert a remediation instruction steps array to form array
 *
 * @param {RemediationInstructionStep[]} steps
 * @returns {Array}
 */
const remediationInstructionStepsToForm = steps => steps.map(step => ({
  ...step,
  operations: remediationInstructionStepOperationsToForm(step.operations),
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
 * Convert a remediation instruction step operations form array to a API compatible operation array
 *
 * @param {RemediationInstructionStepOperation[]} operations
 * @returns {Array}
 */
const formOperationsToRemediationInstructionOperation = operations => operations.map((operation) => {
  const { interval, unit } = operation.time_to_complete;

  return ({
    ...omit(operation, ['key', 'saved']),
    time_to_complete: getSecondsByUnit(interval, unit),
    time_to_complete_unit: unit,
  });
});

/**
 * Convert a remediation instruction steps form array to a API compatible array
 *
 * @param {RemediationInstructionStep[]} steps
 * @returns {Array}
 */
const formStepsToRemediationInstructionSteps = steps => steps.map(step => ({
  ...omit(step, ['key', 'saved']),
  operations: formOperationsToRemediationInstructionOperation(step.operations),
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
