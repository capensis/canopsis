import { isUndefined, omit, pick, cloneDeep } from 'lodash';

import {
  REMEDIATION_INSTRUCTION_APPROVAL_TYPES,
  REMEDIATION_INSTRUCTION_STATUSES,
  REMEDIATION_INSTRUCTION_TYPES,
  TIME_UNITS,
  WORKFLOW_TYPES,
} from '@/constants';

import { uuid } from '@/helpers/uuid';
import { durationToForm } from '@/helpers/date/duration';
import { flattenErrorMap } from '@/helpers/entities/shared/form';

/**
 * @typedef {
 *   'stateinc' | 'statedec' | 'pbhenter' | 'pbhleave' | 'activate' | 'unsnooze'
 * } RemediationInstructionAutoTrigger
 */

/**
 * @typedef {Object} RemediationInstructionStepOperation
 * @property {string} operation_id
 * @property {string} name
 * @property {string} description
 * @property {Array} jobs
 * @property {Duration} time_to_complete
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
 */

/**
 * @typedef {RemediationInstructionStep} RemediationInstructionStepForm
 * @property {RemediationInstructionStepOperationForm[]} operations
 * @property {string} [key]
 */

/**
 * @typedef {Object} RemediationInstructionApprovalUser
 * @property {User} user
 */

/**
 * @typedef {Object} RemediationInstructionApprovalRole
 * @property {Role} role
 */

/**
 * @typedef {Object} RemediationInstructionApprovalRequestedBy
 * @property {string} _id
 * @property {string} name
 * @property {string} display_name
 */

/**
 * @typedef {RemediationInstructionApprovalUser | RemediationInstructionApprovalRole} RemediationInstructionApproval
 * @property {string} comment
 * @property {RemediationInstructionApprovalRequestedBy} requested_by
 * @property {string} [dismissed_comment]
 * @property {RemediationInstructionApprovalRequestedBy} [dismissed_by]
 */

/**
 * @typedef {RemediationInstructionApproval} RemediationInstructionApprovalForm
 * @property {number} type
 */

/**
 * @typedef {RemediationInstructionApproval} RemediationInstructionJob
 * @property {RemediationJob} job
 * @property {boolean} stop_on_fail
 */

/**
 * @typedef {RemediationInstructionJob} RemediationInstructionJobForm
 * @property {string} [key]
 */

/**
 * @typedef {Object} RemediationInstructionManual
 * @property {RemediationInstructionStep[]} steps
 */

/**
 * @typedef {Object} RemediationInstructionAuto
 * @property {RemediationInstructionAutoTrigger[]} triggers
 * @property {number} [priority]
 * @property {RemediationInstructionJob[]} [jobs]
 */

/**
 * @typedef {RemediationInstructionManual | RemediationInstructionAuto} RemediationInstruction
 * @property {number} type
 * @property {number} status
 * @property {string} name
 * @property {boolean} enabled
 * @property {string} description
 * @property {Duration} timeout_after_execution
 * @property {Array} [alarm_pattern]
 * @property {Array} [entity_pattern]
 * @property {string[]} active_on_pbh
 * @property {string[]} disabled_on_pbh
 * @property {RemediationInstructionApproval} approval
 */

/**
 * @typedef {RemediationInstruction} RemediationInstructionForm
 * @property {RemediationInstructionStepForm[]} steps
 * @property {RemediationInstructionJobForm[]} jobs
 * @property {RemediationInstructionApprovalForm} approval
 */

/**
 * @typedef {Object} RemediationInstructionApprovalRequest
 * @property {string} [user]
 * @property {string} [role]
 * @property {string} comment
 */

/**
 * Check instruction status is requested approve
 *
 * @param {RemediationInstruction} instruction
 * @returns {boolean}
 */
export const isApproveRequested = instruction => [
  REMEDIATION_INSTRUCTION_STATUSES.createdAndApproveRequested,
  REMEDIATION_INSTRUCTION_STATUSES.updatedAndApproveRequested,
].includes(instruction.status);

/**
 * Check instruction status is approved
 *
 * @param {RemediationInstruction} instruction
 * @returns {boolean}
 */
export const isInstructionApproved = instruction => instruction.status === REMEDIATION_INSTRUCTION_STATUSES.approved;

/**
 * Check instruction status is dismissed
 *
 * @param {RemediationInstruction} instruction
 * @returns {boolean}
 */
export const isInstructionDismissed = instruction => [
  REMEDIATION_INSTRUCTION_STATUSES.createdAndDismissed,
  REMEDIATION_INSTRUCTION_STATUSES.updatedAndDismissed,
].includes(instruction.status);

/**
 * Check instruction type is auto
 *
 * @param {RemediationInstruction} instruction
 * @returns {boolean}
 */
export const isInstructionAuto = instruction => instruction.type === REMEDIATION_INSTRUCTION_TYPES.auto;

/**
 * Check instruction type is manual
 *
 * @param {RemediationInstruction} instruction
 * @returns {boolean}
 */
export const isInstructionManual = instruction => instruction.type === REMEDIATION_INSTRUCTION_TYPES.manual;

/**
 * Check instruction type is simple manual
 *
 * @param {RemediationInstruction} instruction
 * @returns {boolean}
 */
export const isInstructionSimpleManual = instruction => instruction.type === REMEDIATION_INSTRUCTION_TYPES.simpleManual;

/**
 * Convert a remediation instruction step operation to form
 *
 * @param {RemediationInstructionStepOperation} [operation]
 * @returns {RemediationInstructionStepOperationForm}
 */
export const remediationInstructionStepOperationToForm = (operation = {}) => ({
  name: operation.name || '',
  description: operation.description || '',
  time_to_complete: durationToForm(operation.time_to_complete ?? { value: 1, unit: TIME_UNITS.minute }),
  jobs: operation.jobs ? cloneDeep(operation.jobs) : [],
  key: uuid(),
});

/**
 * Convert a remediation instruction step operation array to form array
 *
 * @param {RemediationInstructionStepOperation[]} [operations = [undefined]]
 * @returns {Array}
 */
const remediationInstructionStepOperationsToForm = (operations = [undefined]) => operations
  .map(remediationInstructionStepOperationToForm);

/**
 * Convert a remediation instruction step to form
 *
 * @param {RemediationInstructionStep} [step]
 * @returns {RemediationInstructionStepForm}
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
 * @param {RemediationInstructionStep[]} [steps = [undefined]]
 * @returns {RemediationInstructionStepForm[]}
 */
const remediationInstructionStepsToForm = (steps = [undefined]) => steps.map(remediationInstructionStepToForm);

/**
 * Convert a remediation instruction approval to form
 *
 * @param {RemediationInstructionApproval} approval
 * @return {RemediationInstructionApprovalForm}
 */
const remediationInstructionApprovalToForm = (approval = {}) => ({
  type: approval?.user?._id
    ? REMEDIATION_INSTRUCTION_APPROVAL_TYPES.user
    : REMEDIATION_INSTRUCTION_APPROVAL_TYPES.role,
  user: approval?.user,
  role: approval?.role,
  comment: approval?.comment ?? '',
});

/**
 * Convert a remediation instruction job to form
 *
 * @param {RemediationInstructionJob} [job = {}]
 * @returns {RemediationInstructionJobForm}
 */
export const remediationInstructionJobToForm = (job = {}) => ({
  job: job.job,
  stop_on_fail: !isUndefined(job.stop_on_fail) ? job.stop_on_fail : WORKFLOW_TYPES.stop,
  key: uuid(),
});

/**
 * Convert a remediation instruction jobs array to form array
 *
 * @param {RemediationInstructionJob[]} [jobs = [undefined]]
 * @returns {RemediationInstructionJobForm[]}
 */
const remediationInstructionJobsToForm = (jobs = [undefined]) => jobs.map(remediationInstructionJobToForm);

/**
 * Convert a remediation instruction object to form object
 *
 * @param {RemediationInstruction} remediationInstruction
 * @returns {RemediationInstructionForm}
 */
export const remediationInstructionToForm = (remediationInstruction = {}) => {
  const form = {
    name: remediationInstruction.name || '',
    priority: remediationInstruction.priority,
    type: !isUndefined(remediationInstruction.type)
      ? remediationInstruction.type
      : REMEDIATION_INSTRUCTION_TYPES.manual,
    enabled: remediationInstruction.enabled ?? true,
    timeout_after_execution: durationToForm(remediationInstruction.timeout_after_execution),
    active_on_pbh: remediationInstruction.active_on_pbh
      ? cloneDeep(remediationInstruction.active_on_pbh)
      : [],
    disabled_on_pbh: remediationInstruction.disabled_on_pbh
      ? cloneDeep(remediationInstruction.disabled_on_pbh)
      : [],
    alarm_pattern: remediationInstruction.alarm_pattern,
    entity_pattern: remediationInstruction.entity_pattern,
    description: remediationInstruction.description || '',
    steps: remediationInstructionStepsToForm(remediationInstruction.steps),
    approval: remediationInstructionApprovalToForm(remediationInstruction.approval),
    jobs: remediationInstructionJobsToForm(remediationInstruction.jobs),
  };

  if (remediationInstruction.triggers) {
    form.triggers = [...remediationInstruction.triggers];
  }

  return form;
};

/**
 * Convert a remediation instruction step operations form array to a API compatible operation array
 *
 * @param {RemediationInstructionJobForm[]} jobs
 * @returns {RemediationInstructionJob[]}
 */
const formJobsToRemediationInstructionJobs = (jobs = []) => jobs.map(job => ({
  ...omit(job, ['key']),

  job: job.job._id,
}));

/**
 * Convert a remediation instruction step operations form array to a API compatible operation array
 *
 * @param {RemediationInstructionStepOperationForm[]} operations
 * @returns {RemediationInstructionStepOperation[]}
 */
const formOperationsToRemediationInstructionOperations = operations => operations.map(operation => ({
  ...omit(operation, ['key']),

  jobs: operation.jobs.map(({ _id }) => _id),
}));

/**
 * Convert a remediation instruction steps form array to a API compatible array
 *
 * @param {RemediationInstructionStepForm[]} steps
 * @returns {RemediationInstructionStep[]}
 */
const formStepsToRemediationInstructionSteps = steps => steps.map(step => ({
  ...omit(step, ['key']),

  operations: formOperationsToRemediationInstructionOperations(step.operations),
}));

/**
 * Convert a remediation instruction approval form
 *
 * @param {RemediationInstructionApprovalForm} approval
 * @returns {RemediationInstructionApprovalRequest | undefined}
 */
const formApprovalToRemediationInstructionApproval = (approval) => {
  if (!approval.comment) {
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
 * @param {RemediationInstructionForm} form
 * @returns {RemediationInstruction}
 */
export const formToRemediationInstruction = (form) => {
  const {
    steps, jobs, priority, triggers, ...instruction
  } = form;

  if (isInstructionManual(form)) {
    instruction.steps = formStepsToRemediationInstructionSteps(steps);
  } else {
    instruction.jobs = formJobsToRemediationInstructionJobs(jobs);

    if (isInstructionAuto(form)) {
      instruction.priority = priority;
      instruction.triggers = triggers;
    }
  }

  return {
    ...instruction,
    alarm_pattern: form.alarm_pattern,
    entity_pattern: form.entity_pattern,
    approval: formApprovalToRemediationInstructionApproval(form.approval),
  };
};

/**
 * Convert error structure to form structure
 *
 * @param {FlattenErrors} errors
 * @param {RemediationInstructionForm} form
 * @return {FlattenErrors}
 */
export const remediationInstructionErrorsToForm = (errors, form) => flattenErrorMap(errors, (errorsObject) => {
  const { jobs, ...errorMessages } = errorsObject;

  if (jobs) {
    errorMessages.jobs = jobs.reduce((acc, messages, index) => {
      const job = form.jobs[index];
      acc[job.key] = messages;

      return acc;
    }, {});
  }

  return errorMessages;
});
