import { REMEDIATION_INSTRUCTION_TYPES } from '@/constants';

/**
 * Checks if the instruction is of type 'manual'.
 *
 * @param {number} type - The instruction type to check.
 * @returns {boolean} True if the instruction type is 'manual', otherwise false.
 */
export const isManualInstructionType = type => type === REMEDIATION_INSTRUCTION_TYPES.manual;

/**
 * Checks if the instruction is of type 'auto'.
 *
 * @param {number} type - The instruction type to check.
 * @returns {boolean} True if the instruction type is 'auto', otherwise false.
 */
export const isAutoInstructionType = type => type === REMEDIATION_INSTRUCTION_TYPES.auto;

/**
 * Checks if the instruction is of type 'simpleManual'.
 *
 * @param {number} type - The instruction type to check.
 * @returns {boolean} True if the instruction type is 'simpleManual', otherwise false.
 */
export const isSimpleManualInstructionType = type => type === REMEDIATION_INSTRUCTION_TYPES.simpleManual;

/**
 * Checks if the instruction is of any manual type ('manual' or 'simpleManual').
 *
 * @param {number} type - The instruction type to check.
 * @returns {boolean} True if the instruction type is 'manual' or 'simpleManual', otherwise false.
 */
export const isAnyManualInstructionType = type => (
  isManualInstructionType(type) || isSimpleManualInstructionType(type)
);
