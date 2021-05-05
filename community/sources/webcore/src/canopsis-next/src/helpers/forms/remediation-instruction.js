import { isUndefined, omit, pick } from 'lodash';

import {
  REMEDIATION_INSTRUCTION_APPROVAL_TYPES,
  REMEDIATION_INSTRUCTION_TYPES,
  TIME_UNITS,
  WORKFLOW_TYPES,
} from '@/constants';

import uuid from '@/helpers/uuid';
import { durationToForm, formToDuration } from '@/helpers/date/duration';

/**
 * @typedef {Object} RemediationInstructionStepOperation
 * @property {string} name
 * @property {string} description
 * @property {Array} jobs
 * @property {DurationForm} time_to_complete
 * @property {string} [key]
 */

/**
 * @typedef {RemediationInstructionStepOperation} RemediationInstructionStepOperationForm
 * @property {string} [key]
 */

/**
 * @typedef {Object} RemediationInstructionStep
 * @property {string} endpoint
 * @property {string} name
 * @property {boolean} stop_on_fail
 * @property {RemediationInstructionStepOperation[]} operations
 * @property {string} [key]
 */

/**
 * @typedef {RemediationInstructionStep} RemediationInstructionStepForm
 * @property {RemediationInstructionStepOperationForm[]} operations
 * @property {string} [key]
 */

/**
 * @typedef {Object} RemediationInstructionApproval
 * @property {string} [user]
 * @property {string} [role]
 * @property {string} comment
 */

/**
 * @typedef {RemediationInstructionApproval} RemediationInstructionApprovalForm
 * @property {boolean} need_approve
 * @property {number} type
 */

/**
 * @typedef {Object} RemediationInstruction
 * @property {number} type
 * @property {string} name
 * @property {boolean} enabled
 * @property {string} description
 * @property {Array} alarm_patterns
 * @property {Array} entity_patterns
 * @property {RemediationInstructionStep[]} steps
 * @property {string[]} active_on_pbh
 * @property {string[]} disabled_on_pbh
 * @property {string[]} jobs
 * @property {RemediationInstructionApproval} approval
 */

/**
 * @typedef {RemediationInstruction} RemediationInstructionForm
 * @property {RemediationInstructionStepForm[]} steps
 * @property {RemediationInstructionApprovalForm} approval
 */

/**
 * Convert a remediation instruction step operation to form
 *
 * @param {RemediationInstructionStepOperation} [operation]
 * @returns {Array}
 */
export const remediationInstructionStepOperationToForm = (operation = {}) => ({
  name: operation.name || '',
  description: operation.description || '',
  time_to_complete: operation.time_to_complete
    ? durationToForm(operation.time_to_complete)
    : { value: 0, unit: TIME_UNITS.minute },
  jobs: operation.jobs || [],
  key: uuid(),
});

/**
 * Convert a remediation instruction step operation array to form array
 *
 * @param {RemediationInstructionStepOperation[]} operations
 * @returns {Array}
 */
const remediationInstructionStepOperationsToForm = (operations = [undefined]) =>
  operations.map(remediationInstructionStepOperationToForm);

/**
 * Convert a remediation instruction step to form
 *
 * @param {RemediationInstructionStep} [step]
 * @returns {Array}
 */
export const remediationInstructionStepToForm = (step = {}) => ({
  endpoint: step.endpoint || '',
  name: step.name || '',
  stop_on_fail: !isUndefined(step.stop_on_fail) ? step.stop_on_fail : WORKFLOW_TYPES.stop,
  operations: remediationInstructionStepOperationsToForm(step.operations),
  key: uuid(),
});

/**
 * Convert a remediation instruction steps array to form array
 *
 * @param {RemediationInstructionStep[]} steps
 * @returns {RemediationInstructionStepForm[]}
 */
const remediationInstructionStepsToForm = (steps = [undefined]) => steps.map(remediationInstructionStepToForm);


const remediationInstructionApprovalToForm = (approval = {}) => ({
  need_approve: !!approval.comment,
  type: approval.user
    ? REMEDIATION_INSTRUCTION_APPROVAL_TYPES.user
    : REMEDIATION_INSTRUCTION_APPROVAL_TYPES.role,
  user: approval.user,
  role: approval.role,
  comment: approval.comment || '',
});

/**
 * Convert a remediation instruction object to form object
 *
 * @param {RemediationInstruction} remediationInstruction
 * @returns {RemediationInstructionForm}
 */
export const remediationInstructionToForm = (remediationInstruction = {}) => ({
  name: remediationInstruction.name || '',
  type: !isUndefined(remediationInstruction.type) ? remediationInstruction.type : REMEDIATION_INSTRUCTION_TYPES.manual,
  enabled: !isUndefined(remediationInstruction.enabled) ? remediationInstruction.enabled : true,
  active_on_pbh: remediationInstruction.active_on_pbh || [],
  disabled_on_pbh: remediationInstruction.disabled_on_pbh || [],
  alarm_patterns: remediationInstruction.alarm_patterns || [],
  entity_patterns: remediationInstruction.entity_patterns || [],
  description: remediationInstruction.description || '',
  steps: remediationInstructionStepsToForm(remediationInstruction.steps),
  approval: remediationInstructionApprovalToForm(remediationInstruction.approval),
  jobs: remediationInstruction.jobs || [],
});


/**
 * Convert a remediation instruction step operations form array to a API compatible operation array
 *
 * @param {RemediationJob[]} jobs
 * @returns {string[]}
 */
const formJobsToRemediationInstructionJobs = (jobs = []) => jobs.map(({ _id }) => _id);


/**
 * Convert a remediation instruction step operations form array to a API compatible operation array
 *
 * @param {RemediationInstructionStepOperationForm[]} operations
 * @returns {RemediationInstructionStepOperation}
 */
const formOperationsToRemediationInstructionOperation = operations => operations.map(operation => ({
  ...omit(operation, ['key']),

  time_to_complete: formToDuration(operation.time_to_complete),
  jobs: formJobsToRemediationInstructionJobs(operation.jobs),
}));

/**
 * Convert a remediation instruction steps form array to a API compatible array
 *
 * @param {RemediationInstructionStepForm[]} steps
 * @returns {RemediationInstructionStep[]}
 */
const formStepsToRemediationInstructionSteps = steps => steps.map(step => ({
  ...omit(step, ['key']),

  operations: formOperationsToRemediationInstructionOperation(step.operations),
}));

/**
 * Convert a remediation instruction approval form
 *
 * @param {RemediationInstructionApprovalForm} approval
 * @returns {RemediationInstructionApproval | undefined}
 */
const formApprovalToRemediationInstructionApproval = (approval) => {
  if (!approval.need_approve) {
    return undefined;
  }

  const data = pick(approval, ['comment']);

  if (approval.type === REMEDIATION_INSTRUCTION_APPROVAL_TYPES.role) {
    data.role = approval.role._id;
  } else {
    data.user = approval.user._id;
  }

  return data;
};

/**
 * Convert a remediation instruction form object to a API compatible object
 *
 * @param {Object} form
 * @returns {Object}
 */
export const formToRemediationInstruction = (form) => {
  const {
    steps, jobs, priority, ...instruction
  } = form;

  if (form.type === REMEDIATION_INSTRUCTION_TYPES.manual) {
    instruction.steps = formStepsToRemediationInstructionSteps(steps);
  } else {
    instruction.priority = priority;
    instruction.jobs = formJobsToRemediationInstructionJobs(jobs);
  }

  return {
    ...instruction,
    alarm_patterns: form.alarm_patterns.length ? form.alarm_patterns : undefined,
    entity_patterns: form.entity_patterns.length ? form.entity_patterns : undefined,
    approval: formApprovalToRemediationInstructionApproval(form.approval),
  };
};
