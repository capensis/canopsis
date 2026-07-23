import { REMEDIATION_INSTRUCTION_STATUSES, REMEDIATION_INSTRUCTION_TYPES } from '@/constants';

/**
 * @typedef {
 *   'stateinc' | 'statedec' | 'pbhenter' | 'pbhleave' | 'activate' | 'unsnooze'
 * } RemediationInstructionAutoTrigger
 */

/**
 * @typedef {Object} RemediationInstructionApprovalRequestedBy
 * @property {string} _id
 * @property {string} name
 * @property {string} display_name
 */

/**
 * @typedef {Object} RemediationInstructionApproval
 * @property {string} comment
 * @property {RemediationInstructionApprovalRequestedBy} requested_by
 * @property {string} [dismissed_comment]
 * @property {RemediationInstructionApprovalRequestedBy} [dismissed_by]
 */

/**
 * @typedef {Object} RemediationInstruction
 * @property {number} type
 * @property {number} status
 * @property {string} name
 * @property {boolean} enabled
 * @property {string} description
 * @property {Array} [alarm_pattern]
 * @property {Array} [entity_pattern]
 * @property {string[]} active_on_pbh
 * @property {string[]} disabled_on_pbh
 * @property {RemediationInstructionApproval} approval
 */

/**
 * Check instruction status is requested approve
 *
 * @param {RemediationInstruction} instruction
 * @returns {boolean}
 */
export const isApproveRequested = (instruction) =>
  [
    REMEDIATION_INSTRUCTION_STATUSES.createdAndApproveRequested,
    REMEDIATION_INSTRUCTION_STATUSES.updatedAndApproveRequested,
  ].includes(instruction.status);

/**
 * Check instruction status is approved
 *
 * @param {RemediationInstruction} instruction
 * @returns {boolean}
 */
export const isInstructionApproved = (instruction) =>
  instruction.status === REMEDIATION_INSTRUCTION_STATUSES.approved;

/**
 * Check instruction status is dismissed
 *
 * @param {RemediationInstruction} instruction
 * @returns {boolean}
 */
export const isInstructionDismissed = (instruction) =>
  [
    REMEDIATION_INSTRUCTION_STATUSES.createdAndDismissed,
    REMEDIATION_INSTRUCTION_STATUSES.updatedAndDismissed,
  ].includes(instruction.status);

/**
 * Check instruction type is auto
 *
 * @param {number} type
 * @returns {boolean}
 */
export const isInstructionTypeAuto = (type) => type === REMEDIATION_INSTRUCTION_TYPES.auto;

/**
 * Check instruction type is manual
 *
 * @param {number} type
 * @returns {boolean}
 */
export const isInstructionTypeManual = (type) => type === REMEDIATION_INSTRUCTION_TYPES.manual;

/**
 * Check instruction type is simple manual
 *
 * @param {number} type
 * @returns {boolean}
 */
export const isInstructionTypeSimpleManual = (type) =>
  type === REMEDIATION_INSTRUCTION_TYPES.simpleManual;

/**
 * Check instruction type is any manual
 *
 * @param {number} type
 * @returns {boolean}
 */
export const isInstructionTypeAnySimpleManual = (type) =>
  isInstructionTypeManual(type) || isInstructionTypeSimpleManual(type);
