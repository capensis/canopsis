import { isString, kebabCase } from 'lodash';

/**
 * @typedef {Object} ThemeMainColors
 * @property {string} primary
 * @property {string} secondary
 * @property {string} accent
 * @property {string} error
 * @property {string} error_icons
 * @property {string} info
 * @property {string} info_icons
 * @property {string} success
 * @property {string} success_icons
 * @property {string} warning
 * @property {string} warning_icons
 * @property {string} background
 */

/**
 * @typedef {Object} ThemeTableColors
 * @property {string} background
 * @property {string} row_color
 * @property {string} shift_row_color
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

/**
 * Converts main theme colors to Vuetify-compatible variables.
 *
 * @param {Object} [main={}] - The main theme colors object.
 * @param {string} [main.error_icons] - The color for error icons.
 * @param {string} [main.error] - The background color for errors.
 * @param {string} [main.warning_icons] - The color for warning icons.
 * @param {string} [main.warning] - The background color for warnings.
 * @param {string} [main.success_icons] - The color for success icons.
 * @param {string} [main.success] - The background color for success.
 * @param {string} [main.info_icons] - The color for info icons.
 * @param {string} [main.info] - The background color for info.
 * @returns {Object} An object with Vuetify-compatible variable names.
 */
export const convertMainToVuetifyVariables = (main = {}) => ({
  ...main,
  error: main.error_icons,
  error_background: main.error,

  warning: main.warning_icons,
  warning_background: main.warning,

  success: main.success_icons,
  success_background: main.success,

  info: main.info_icons,
  info_background: main.info,
});
