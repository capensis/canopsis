import tinycolor from 'tinycolor2';

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
 * Get darken color
 *
 * @param {Element} element
 * @param {string} property
 */
export const getCSSVariableColor = (element, property) => getComputedStyle(element)
  .getPropertyValue(property);
