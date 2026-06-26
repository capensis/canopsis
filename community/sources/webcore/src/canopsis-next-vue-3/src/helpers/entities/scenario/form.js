import { DEPRECATED_TRIGGERS } from '@/constants';

/**
 * Check trigger is deprecated
 *
 * @param {string} type
 * @returns {boolean}
 */
export const isDeprecatedTrigger = type => DEPRECATED_TRIGGERS.includes(type);
