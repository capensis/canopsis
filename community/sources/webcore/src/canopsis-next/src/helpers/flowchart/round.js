/**
 * Round value by step
 *
 * @param {number} value
 * @param {number} [step = 1]
 * @returns {number}
 */
export const roundByStep = (value, step) => Math.round(value / step) * step;
