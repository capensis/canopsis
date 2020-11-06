import { FIRST_LETTER_ALPHABET_CHAR_CODE } from '@/constants';

/**
 * Get char by index
 *
 * @param {number} index
 * @return {string}
 */
export const getCharByIndex = index => String.fromCharCode(FIRST_LETTER_ALPHABET_CHAR_CODE + index);
