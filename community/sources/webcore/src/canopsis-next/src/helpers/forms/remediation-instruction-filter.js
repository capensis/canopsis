import { REMEDIATION_INSTRUCTION_TYPES } from '@/constants';

import { enabledToForm } from './shared/common';

/**
 * @typedef {Object} RemediationInstructionFilterInstruction
 * @property {string} _id
 * @property {string} name
 * @property {number} type
 */

/**
 * @typedef {Object} RemediationInstructionFilterForm
 * @property {boolean} with
 * @property {boolean} all
 * @property {boolean} all
 * @property {boolean} auto
 * @property {boolean} manual
 * @property {boolean} [has_running]
 * @property {RemediationInstructionFilterInstruction[]} instructions
 */

/**
 * @typedef {RemediationInstructionFilterForm} RemediationInstructionFilter
 * @property {string} _id
 * @property {boolean} locked
 */

/**
 * Convert remediation instruction filter to form
 *
 * @param {RemediationInstructionFilter | {}} [filter = {}]
 * @return {RemediationInstructionFilterForm}
 */
export const remediationInstructionFilterToForm = (filter = {}) => ({
  with: enabledToForm(filter.with),
  all: !!filter.all,
  auto: !!filter.auto,
  manual: !!filter.manual,
  has_running: filter.has_running,
  instructions: filter.instructions ? [...filter.instructions] : [],
});

/**
 * Is remediation instruction intersects with remediation instruction filter by type
 *
 * @param {RemediationInstructionFilterForm | {}} [filter = {}]
 * @param {RemediationInstruction | {}} [instruction = {}]
 * @returns {boolean}
 */
export const isRemediationInstructionIntersectsWithFilterByType = (filter = {}, instruction = {}) => {
  const isAuto = instruction.type === REMEDIATION_INSTRUCTION_TYPES.auto;

  return (filter.auto && isAuto) || (filter.manual && !isAuto);
};
