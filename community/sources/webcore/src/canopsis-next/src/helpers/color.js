import tinycolor from 'tinycolor2';

import { BACKGROUND_AND_ICONS_COLORS_DIFF } from '@/constants';

/**
 * @typedef {Object} ReadableColorOptions
 * @property {'AA' | 'AAA'} level
 * @property {'small' | 'large'} size
 */

/**
 * Check color is readable
 *
 * @param {string} firstColor
 * @param {string} secondColor
 * @param {ReadableColorOptions} [options = {}]
 * @return {boolean}
 */
export const isReadableColor = (firstColor, secondColor, options = {}) => tinycolor.isReadable(
  firstColor,
  secondColor,
  options,
);

/**
 * Get most readable text color ('white' or 'black')
 *
 * @param {string} color
 * @param {ReadableColorOptions} [options = {}]
 */
export const getMostReadableTextColor = (color, options = {}) => {
  if (!color) {
    return 'black';
  }

  const isWhiteReadable = isReadableColor(color, 'white', options);

  return isWhiteReadable ? 'white' : 'black';
};

/**
 * Convert color to rgb
 *
 * @param {string|Object} color
 * @return {string}
 */
export const colorToRgb = color => tinycolor(color).toRgbString();

/**
 * Convert color to rgba with alpha
 *
 * @param {string|Object} color
 * @param {number} alpha
 * @return {string}
 */
export const colorToRgba = (color, alpha = 1.0) => tinycolor(color)
  .setAlpha(alpha)
  .toRgbString();

/**
 * Convert color to hex
 *
 * @param {string|Object} color
 * @return {string}
 */
export const colorToHex = color => tinycolor(color).toHexString();

/**
 * Check color is valid
 *
 * @param {string|Object} color
 * @return {boolean}
 */
export const isValidColor = color => tinycolor(color).isValid();

/**
 * Check color is dark
 *
 * @param {string|Object} color
 * @return {boolean}
 */
export const isDarkColor = color => tinycolor(color).isDark();

/**
 * Get darken color
 *
 * @param {string} color
 * @param {number} amount
 */
export const getDarkenColor = (color, amount) => tinycolor(color)
  .darken(amount)
  .toString();

/**
 * Get lighten color
 *
 * @param {string} color
 * @param {number} amount
 */
export const getLightenColor = (color, amount) => tinycolor(color)
  .lighten(amount)
  .toString();

/**
 * Adjusts the color by either lightening or darkening it based on its current brightness.
 *
 * @param {string|Object} color - The color to be adjusted. It can be a string or an object compatible with tinycolor.
 * @param {number} [amount=BACKGROUND_AND_ICONS_COLORS_DIFF] - The amount by which to lighten or darken the color.
 * @returns {string} - The adjusted color as a string.
 */
export const getLightenOrDarkenColor = (color, amount = BACKGROUND_AND_ICONS_COLORS_DIFF) => (
  isDarkColor(color) ? getLightenColor(color, amount) : getDarkenColor(color, amount)
);

/**
 * Check property is css variable
 *
 * @param {string} property
 * @returns {boolean}
 */
export const isCSSVariable = property => /^var\(.+\)$/.test(property);

/**
 * Get css variable name
 *
 * @param {string} property
 * @returns {string}
 */
export const getCSSVariableName = property => property.match(/^var\((.+)\)$/)[1];

/**
 * Get color from css variable
 *
 * @param {Element} element
 * @param {string} property
 */
export const getCSSVariableColor = (element, property) => getComputedStyle(element)
  .getPropertyValue(property);
