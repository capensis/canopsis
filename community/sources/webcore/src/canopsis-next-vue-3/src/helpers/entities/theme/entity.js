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
 * @typedef {Object} ThemeWidgetColors
 * @property {string} title_background
 * @property {string} content_background
 * @property {boolean} border
 * @property {string} border_color
 * @property {number} border_radius
 */

/**
 * @typedef {Object} ThemeStatesColors
 * @property {string} ok
 * @property {string} minor
 * @property {string} major
 * @property {string} critical
 */

/**
 * @typedef {Object} ThemeIconsColors
 * @property {string} ack
 * @property {string} canceled
 * @property {string} children
 * @property {string} parent
 * @property {string} comment
 * @property {string} snooze
 * @property {string} tickets
 * @property {string} cancel_delay
 */

/**
 * @typedef {Object} ThemeColors
 * @property {ThemeMainColors} main
 * @property {ThemeTableColors} table
 * @property {ThemeWidgetColors} widget
 * @property {ThemeStatesColors} state
 * @property {ThemeIconsColors} icons
 */

/**
 * @typedef {Object} Theme
 * @property {String} name
 * @property {ThemeColors} colors
 * @property {number} font_size
 */

const themeObjectColorsToCSSVariables = (colors, prefix = '') =>
  Object.entries(colors).reduce((acc, [key, value]) => {
    if (!value) {
      return acc;
    }

    if (isString(value) || typeof value === 'number') {
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
export const themePropertiesToCSSVariables = (colors = {}) =>
  themeObjectColorsToCSSVariables(colors);

/**
 * Extracts the theme variable key from a CSS variable string.
 *
 * @param {string} [cssVariable=''] - The CSS variable string to process.
 * @returns {string} The extracted theme variable key.
 */
export const getThemeVariableKeyFromCssVariable = (cssVariable = '') =>
  cssVariable.replace('var(--v-', '').replace('-base)', '');
