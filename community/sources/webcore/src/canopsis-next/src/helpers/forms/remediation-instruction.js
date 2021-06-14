import { isUndefined, omit } from 'lodash';

import uuid from '@/helpers/uuid';
import { durationToForm, formToDuration } from '@/helpers/date/duration';
import { generateRemediationInstructionStep } from '@/helpers/entities';

/**
 * Convert a remediation instruction step operation array to form array
 *
 * @param {RemediationInstructionStepOperation[]} operations
 * @returns {Array}
 */
const remediationInstructionStepOperationsToForm = operations => operations.map(operation => ({
  ...operation,

  time_to_complete: durationToForm(operation.time_to_complete),
  jobs: operation.jobs || [],
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
 * @property {Array} alarm_patterns
 * @property {Array} entity_patterns
 * @property {RemediationInstructionStep[]} steps
 * @param {Object} remediationInstruction
 * @returns {RemediationInstruction}
 */
export const remediationInstructionToForm = (remediationInstruction = {}) => ({
  name: remediationInstruction.name || '',
  enabled: !isUndefined(remediationInstruction.enabled) ? remediationInstruction.enabled : true,
  active_on_pbh: remediationInstruction.active_on_pbh || [],
  disabled_on_pbh: remediationInstruction.disabled_on_pbh || [],
  alarm_patterns: remediationInstruction.alarm_patterns || [],
  entity_patterns: remediationInstruction.entity_patterns || [],
  description: remediationInstruction.description || '',
  steps: remediationInstruction.steps
    ? remediationInstructionStepsToForm(remediationInstruction.steps)
    : [generateRemediationInstructionStep()],
});


/**
 * Convert a remediation instruction step operations form array to a API compatible operation array
 *
 * @param {RemediationInstructionStepOperation[]} operations
 * @returns {Array}
 */
const formOperationsToRemediationInstructionOperation = operations => operations.map(operation => ({
  ...omit(operation, ['key']),

  time_to_complete: formToDuration(operation.time_to_complete),
  jobs: operation.jobs.map(({ _id }) => _id),
}));

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

  alarm_patterns: form.alarm_patterns && form.alarm_patterns.length ? form.alarm_patterns : undefined,
  entity_patterns: form.entity_patterns && form.entity_patterns.length ? form.entity_patterns : undefined,
  steps: formStepsToRemediationInstructionSteps(form.steps),
});
