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
 * @property {ThemeStatesColors} states
 */

/**
 * @typedef {Object} Theme
 * @property {String} name
 * @property {ThemeColors} colors
 */

/**
 * @typedef {ThemeMainColors} ThemeMainColorsForm
 */

/**
 * @typedef {ThemeTableColors} ThemeTableColorsForm
 */

/**
 * @typedef {ThemeStatesColors} ThemeStatesColorsForm
 */

/**
 * @typedef {User} ThemeColorsForm
 * @property {ThemeMainColorsForm} main
 * @property {ThemeTableColorsForm} table
 * @property {ThemeStatesColorsForm} states
 */

/**
 * @typedef {User} ThemeForm
 * @property {ThemeColorsForm} colors
 */

import { COLORS } from '@/config';

/**
 * Convert theme main colors to form object
 *
 * @param {ThemeMainColors} [main = {}]
 * @returns {ThemeMainColorsForm}
 */
export const themeMainColorsToForm = (main = {}) => ({
  primary: main.primary ?? COLORS.primary,
  secondary: main.secondary ?? COLORS.secondary,
  accent: main.accent ?? COLORS.accent,
  error: main.error ?? COLORS.error,
  info: main.info ?? COLORS.info,
  success: main.success ?? COLORS.success,
  warning: main.warning ?? COLORS.warning,
  background: main.background ?? COLORS.background,
});

/**
 * Convert theme main colors to form object
 *
 * @param {ThemeTableColors} [table = {}]
 * @returns {ThemeTableColorsForm}
 */
export const themeTableColorsToForm = (table = {}) => ({
  background: table.background ?? COLORS.table.background,
  active_color: table.active_color ?? COLORS.table.activeColor,
  row_color: table.row_color ?? COLORS.table.rowColor,
  shift_row_color: table.shift_row_color,
  hover_row_color: table.hover_row_color ?? COLORS.table.hoverRowColor,
});

/**
 * Convert theme main colors to form object
 *
 * @param {ThemeStatesColors} [states = {}]
 * @returns {ThemeStatesColorsForm}
 */
export const themeStatesColorsToForm = (states = {}) => ({
  ok: states.ok ?? COLORS.state.ok,
  minor: states.minor ?? COLORS.state.minor,
  major: states.major ?? COLORS.state.major,
  critical: states.critical ?? COLORS.state.critical,
});

/**
 * Convert theme colors to form object
 *
 * @param {ThemeColors} [colors = {}]
 * @returns {ThemeColorsForm}
 */
export const themeColorsToForm = (colors = {}) => ({
  main: themeMainColorsToForm(colors.main),
  table: themeTableColorsToForm(colors.table),
  states: themeStatesColorsToForm(colors.states),
});

/**
 * Convert theme to form object
 *
 * @param {Theme} [theme = {}]
 * @returns {ThemeForm}
 */
export const themeToForm = (theme = {}) => ({
  name: theme.name ?? '',
  colors: themeColorsToForm(theme.colors),
});
