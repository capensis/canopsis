const FIRST_LETTER_ALPHABET_CHAR_CODE = 97;

/**
 * Get letter by index
 *
 * @param {number} index
 * @return {string}
 */
export const getLetterByIndex = (index = 0) => String.fromCharCode(FIRST_LETTER_ALPHABET_CHAR_CODE + index);
