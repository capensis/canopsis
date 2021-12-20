import { isNil } from 'lodash';

const FIRST_LETTER_ALPHABET_CHAR_CODE = 97;

/**
 * Get letter by index
 *
 * @param {number} index
 * @return {string}
 */
export const getLetterByIndex = (index = 0) => String.fromCharCode(FIRST_LETTER_ALPHABET_CHAR_CODE + index);

/**
 * Convert number to string with fixed point
 *
 * @param {number} value
 * @param {number} [digits = 3]
 * @return {string}
 */
export const convertNumberToFixedString = (value, digits = 3) => value && Number(value).toFixed(digits);

/**
 * Convert number to rounded percent string
 *
 * @param {number} value - Numeric value to format
 * @param {number} [precision = 3] - Number of floating digit to keep
 *
 * @returns {string}
 */
export const convertNumberToRoundedPercentString = (value, precision = 3) => {
  if (isNil(value)) {
    return '';
  }

  const tmp = 10 ** precision;

  const filteredValue = 100 * value;

  const roundedValue = Math.round(filteredValue * tmp) / tmp;

  return `${roundedValue}%`;
};
