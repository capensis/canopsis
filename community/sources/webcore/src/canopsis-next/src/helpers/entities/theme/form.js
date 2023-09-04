import { COLORS } from '@/config';

/**
 * @typedef {Object} ThemeEnabledColor
 * @property {boolean} enabled
 * @property {color} string
 */

/**
 * @typedef {ThemeMainColors} ThemeMainColorsForm
 */

/**
 * @typedef {ThemeTableColors} ThemeTableColorsForm
 * @property {ThemeEnabledColor} shift_row_color
 * @property {ThemeEnabledColor} hover_row_color
 */

/**
 * @typedef {ThemeStatesColors} ThemeStatesColorsForm
 */

/**
 * @typedef {User} ThemeColorsForm
 * @property {ThemeMainColorsForm} main
 * @property {ThemeTableColorsForm} table
 * @property {ThemeStatesColorsForm} state
 */

/**
 * @typedef {User} ThemeForm
 * @property {ThemeColorsForm} colors
 */

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
  shift_row_color: {
    enabled: table.shift_row_color ?? false,
    color: table.shift_row_color ?? COLORS.table.shiftRowColor,
  },
  hover_row_color: {
    enabled: table.hover_row_color ?? true,
    color: table.hover_row_color ?? COLORS.table.hoverRowColor,
  },
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
  state: themeStatesColorsToForm(colors.state),
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
