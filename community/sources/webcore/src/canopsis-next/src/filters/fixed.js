/**
 * Convert number to string with fixed point
 *
 * @param {number} value
 * @param {number} digits
 * @return {string}
 */
export default function (value, digits = 3) {
  return value && Number(value).toFixed(digits);
}
