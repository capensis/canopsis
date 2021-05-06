/**
 * @typedef {Object} RemediationInstructionFilterInstruction
 * @property {string} _id
 * @property {string} name
 */

/**
 * @typedef {Object} RemediationInstructionFilterForm
 * @property {boolean} with
 * @property {boolean} all
 * @property {boolean} all
 * @property {boolean} auto
 * @property {boolean} manual
 * @property {RemediationInstructionFilterInstruction[]} instructions
 */

/**
 * @typedef {RemediationInstructionFilterForm} RemediationInstructionFilter
 * @property {string} _id
 * @property {boolean} locked
 */

import { isUndefined } from 'lodash';

/**
 * Convert remediation instruction filter to form
 *
 * @param {RemediationInstructionFilter | {}} [filter = {}]
 * @return {RemediationInstructionFilterForm}
 */
export const remediationInstructionFilterToForm = (filter = {}) => ({
  with: isUndefined(filter.with) ? true : filter.with,
  all: !!filter.all,
  auto: !!filter.auto,
  manual: !!filter.manual,
  instructions: filter.instructions ? [...filter.instructions] : [],
});
