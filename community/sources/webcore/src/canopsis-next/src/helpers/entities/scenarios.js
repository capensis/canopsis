import { DEPRECATED_TRIGGERS } from '@/constants';

/**
 * Check trigger is deprecated
 *
 * @param {string} trigger
 * @returns {boolean}
 */
export const isDeprecatedTrigger = trigger => DEPRECATED_TRIGGERS.includes(trigger);
