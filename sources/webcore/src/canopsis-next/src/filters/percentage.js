/**
 *
 * @param {Number} value - Numeric value to format
 * @param {Number} precision - Number of floating digit to keep
 *
 * @returns {String}
 */
export default function (value, precision = 3) {
  const tmp = 10 ** precision;

  const filteredValue = 100 * value;

  const roundedValue = Math.round(filteredValue * tmp) / tmp;

  return `${roundedValue}%`;
}
