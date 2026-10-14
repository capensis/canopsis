import { isNumber } from 'lodash';

/**
 * Clamp a number to be at least 0. Non-number values are returned unchanged.
 *
 * @param {*} value
 * @return {*}
 */
export const normalizeNonNegativeNumber = value => (isNumber(value) ? Math.max(value, 0) : value);

/**
 * Clamp a number to be at least 1. Non-number values are returned unchanged.
 *
 * @param {*} value
 * @return {*}
 */
export const normalizePositiveNumber = value => (isNumber(value) ? Math.max(value, 1) : value);
