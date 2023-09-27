import { COLORS } from '@/config';

import { THEME_FONT_SIZES } from '@/constants/theme';

/**
 * @typedef {Object} ThemeEnabledColor
 * @property {boolean} enabled
 * @property {string} color
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
 * @typedef {ThemeColors} ThemeColorsForm
 * @property {ThemeMainColorsForm} main
 * @property {ThemeTableColorsForm} table
 * @property {ThemeStatesColorsForm} state
 */

/**
 * @typedef {Theme} ThemeForm
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
  active_color: main.active_color ?? COLORS.activeColor,
});

/**
 * Convert theme main colors to form object
 *
 * @param {ThemeTableColors} [table = {}]
 * @returns {ThemeTableColorsForm}
 */
export const themeTableColorsToForm = (table = {}) => ({
  background: table.background ?? COLORS.table.background,
  row_color: table.row_color ?? COLORS.table.rowColor,
  shift_row_color: {
    enabled: !!table.shift_row_color,
    color: table.shift_row_color ?? COLORS.table.shiftRowColor,
  },
  hover_row_color: {
    enabled: !!table.hover_row_color,
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
  font_size: theme.font_size ?? THEME_FONT_SIZES.medium,
  colors: themeColorsToForm(theme.colors),
});

/**
 * Convert theme enabled color to string
 *
 * @param {ThemeEnabledColor} value
 * @returns {string}
 */
const themeEnabledColorToString = value => (value.enabled ? value.color : null);

/**
 * Convert theme table colors form to API compatible object
 *
 * @param {ThemeTableColorsForm} table
 * @returns {ThemeTableColors}
 */
const formTableColorsToTheme = table => ({
  ...table,
  shift_row_color: themeEnabledColorToString(table.shift_row_color),
  hover_row_color: themeEnabledColorToString(table.hover_row_color),
});

/**
 * Convert theme colors form to API compatible object
 *
 * @param {ThemeColorsForm} colors
 * @returns {ThemeColors}
 */
const formColorsToTheme = colors => ({
  ...colors,
  table: formTableColorsToTheme(colors.table),
});

/**
 * Convert theme form to API compatible object
 *
 * @param {ThemeForm} form
 * @returns {Theme}
 */
export const formToTheme = form => ({
  ...form,
  colors: formColorsToTheme(form.colors),
});
