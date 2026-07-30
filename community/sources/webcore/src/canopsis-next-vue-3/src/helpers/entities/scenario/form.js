import { DEPRECATED_TRIGGERS } from '@/constants';

/**
 * @typedef {Object} Trigger
 * @property {string} type
 * @property {number} [threshold]
 */

/**
 * @typedef {Scenario} ScenarioForm
 * @property {ActionForm[]} actions
 */

/**
 * Check trigger is deprecated
 *
 * @param {string} type
 * @returns {boolean}
 */
export const isDeprecatedTrigger = type => DEPRECATED_TRIGGERS.includes(type);
