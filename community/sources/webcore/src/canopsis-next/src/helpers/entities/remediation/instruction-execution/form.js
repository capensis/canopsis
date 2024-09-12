import { INSTRUCTION_EXECUTION_ICONS, REMEDIATION_INSTRUCTION_EXECUTION_STATUSES } from '@/constants';

/**
 * Check type is manual in progress
 *
 * @param {number} icon
 * @returns {boolean}
 */
export const isInstructionExecutionIconManualInProgress = icon => [
  INSTRUCTION_EXECUTION_ICONS.manualInProgress,
  INSTRUCTION_EXECUTION_ICONS.manualFailedWithInProgress,
  INSTRUCTION_EXECUTION_ICONS.manualSuccessfulWithInProgress,
].includes(icon);

/**
 * Check type is auto in progress
 *
 * @param {number} icon
 * @returns {boolean}
 */
export const isInstructionExecutionIconAutoInProgress = icon => [
  INSTRUCTION_EXECUTION_ICONS.autoInProgress,
  INSTRUCTION_EXECUTION_ICONS.autoFailedWithInProgress,
  INSTRUCTION_EXECUTION_ICONS.autoSuccessfulWithInProgress,
].includes(icon);

/**
 * Check type is in progress
 *
 * @param {number} icon
 * @returns {boolean}
 */
export const isInstructionExecutionIconInProgress = icon => isInstructionExecutionIconManualInProgress(icon)
  || isInstructionExecutionIconAutoInProgress(icon);

/**
 * Check type is failed
 *
 * @param {number} icon
 * @returns {boolean}
 */
export const isInstructionExecutionIconFailed = icon => [
  INSTRUCTION_EXECUTION_ICONS.autoFailed,
  INSTRUCTION_EXECUTION_ICONS.manualFailed,
  INSTRUCTION_EXECUTION_ICONS.manualFailedWithInProgress,
  INSTRUCTION_EXECUTION_ICONS.autoFailedWithInProgress,
  INSTRUCTION_EXECUTION_ICONS.autoFailedWithManualAvailable,
  INSTRUCTION_EXECUTION_ICONS.manualFailedWithManualAvailable,
].includes(icon);

/**
 * Check type is success
 *
 * @param {number} icon
 * @returns {boolean}
 */
export const isInstructionExecutionIconSuccess = icon => [
  INSTRUCTION_EXECUTION_ICONS.autoSuccessful,
  INSTRUCTION_EXECUTION_ICONS.manualSuccessful,
  INSTRUCTION_EXECUTION_ICONS.manualSuccessfulWithInProgress,
  INSTRUCTION_EXECUTION_ICONS.autoSuccessfulWithInProgress,
  INSTRUCTION_EXECUTION_ICONS.autoSuccessfulWithManualAvailable,
  INSTRUCTION_EXECUTION_ICONS.manualSuccessfulWithManualAvailable,
].includes(icon);

/**
 * Check type is manual
 *
 * @param {number} icon
 * @returns {boolean}
 */
export const isInstructionExecutionManual = icon => [
  INSTRUCTION_EXECUTION_ICONS.manualInProgress,
  INSTRUCTION_EXECUTION_ICONS.manualFailed,
  INSTRUCTION_EXECUTION_ICONS.manualFailedWithInProgress,
  INSTRUCTION_EXECUTION_ICONS.manualFailedWithManualAvailable,
  INSTRUCTION_EXECUTION_ICONS.manualAvailable,
  INSTRUCTION_EXECUTION_ICONS.manualSuccessful,
  INSTRUCTION_EXECUTION_ICONS.manualSuccessfulWithInProgress,
  INSTRUCTION_EXECUTION_ICONS.manualSuccessfulWithManualAvailable,
].includes(icon);

/**
 * Check instruction execution status is completed
 *
 * @param {number} status
 * @returns {boolean}
 */
export const isInstructionExecutionCompleted = ({
  status,
}) => status === REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.completed;

/**
 * Check instruction execution status is failed
 *
 * @param {number} status
 * @returns {boolean}
 */
export const isInstructionExecutionFailed = ({
  status,
} = {}) => status === REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.failed;

/**
 * Check instruction execution status is aborted
 *
 * @param {number} status
 * @returns {boolean}
 */
export const isInstructionExecutionAborted = ({
  status,
} = {}) => status === REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.aborted;

/**
 * Check instruction execution status is paused
 *
 * @param {number} status
 * @returns {boolean}
 */
export const isInstructionExecutionPaused = ({
  status,
} = {}) => status === REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.paused;

/**
 * Check instruction execution status is running
 *
 * @param {number} status
 * @returns {boolean}
 */
export const isInstructionExecutionRunning = ({
  status,
} = {}) => status === REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.running;
