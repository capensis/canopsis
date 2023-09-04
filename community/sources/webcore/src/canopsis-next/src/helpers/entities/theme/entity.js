/**
 * @typedef {Object} ThemeMainColors
 * @property {string} primary
 * @property {string} secondary
 * @property {string} accent
 * @property {string} error
 * @property {string} info
 * @property {string} success
 * @property {string} warning
 * @property {string} background
 */

/**
 * @typedef {Object} ThemeTableColors
 * @property {string} background
 * @property {string} active_color
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
 */

/**
 * @typedef {Object & ThemeMainColors} ThemeCSSVariables
 */

import { kebabCase } from 'lodash';

const themeObjectColorsToCSSVariables = (colors, prefix = '') => Object.entries(colors)
  .reduce((acc, [key, value]) => {
    acc[`${prefix}${kebabCase(key)}`] = value;

    return acc;
  }, {});

/**
 * Convert theme to form object
 *
 * @param {ThemeColors} [colors = {}]
 * @returns {ThemeCSSVariables}
 */
export const themeColorsToCSSVariables = (colors = {}) => ({
  ...colors.main,
  ...themeObjectColorsToCSSVariables(colors.state, 'state-'),
  ...themeObjectColorsToCSSVariables(colors.table, 'table-'),
});
