import { isString, kebabCase } from 'lodash';

/**
 * @typedef {Object} ThemeMainColors
 * @property {string} primary
 * @property {string} secondary
 * @property {string} accent
 * @property {string} error
 * @property {string} error_background
 * @property {string} info
 * @property {string} info_background
 * @property {string} success
 * @property {string} success_background
 * @property {string} warning
 * @property {string} warning_background
 * @property {string} background
 * @property {string} active_color
 */

/**
 * @typedef {Object} ThemeTableColors
 * @property {string} background
 * @property {string} row_color
 * @property {boolean} shift_row
 * @property {string} shift_row_color
 * @property {boolean} hover_row
 * @property {string} hover_row_color
 */

/**
 * @typedef {Object} ThemeStatesColors
 * @property {string} ok
 * @property {string} minor
 * @property {string} major
 * @property {string} critical
 */

/**
 * @typedef {Object} ThemeColors
 * @property {ThemeMainColors} main
 * @property {ThemeTableColors} table
 * @property {ThemeStatesColors} state
 */

/**
 * @typedef {Object} Theme
 * @property {String} name
 * @property {ThemeColors} colors
 * @property {number} font_size
 */

/**
 * Convert object deep object to flat object variables
 *
 * @param {Object} colors
 * @param {string} prefix
 * @returns {Object}
 */
const themeObjectColorsToCSSVariables = (colors, prefix = '') => Object.entries(colors)
  .reduce((acc, [key, value]) => {
    if (!value) {
      return acc;
    }

    if (isString(value)) {
      acc[`${prefix}${kebabCase(key)}`] = value;

      return acc;
    }

    return {
      ...acc,
      ...themeObjectColorsToCSSVariables(value, `${key}-`),
    };
  }, {});

/**
 * Convert theme to form object
 *
 * @param {Object} [colors = {}]
 * @returns {Object}
 */
export const themePropertiesToCSSVariables = (colors = {}) => themeObjectColorsToCSSVariables(colors);

/**
 * Extracts the theme variable key from a CSS variable string.
 *
 * @param {string} [cssVariable=''] - The CSS variable string to process.
 * @returns {string} The extracted theme variable key.
 */
export const getThemeVariableKeyFromCssVariable = (cssVariable = '') => (
  cssVariable.replace('var(--v-', '').replace('-base)', '')
);
