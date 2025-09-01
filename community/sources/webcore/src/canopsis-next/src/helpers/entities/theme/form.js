import { COLORS } from '@/config';
import { THEME_FONT_SIZES } from '@/constants';

import { isDarkColor, getLightenOrDarkenColor } from '@/helpers/color';

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
 * Convert main type colors to a standardized set of colors with backgrounds.
 *
 * @param {Object} main - The main color object containing type colors.
 * @param {string} [main.error] - The error color.
 * @param {string} [main.error_background] - The error background color.
 * @param {string} [main.warning] - The warning color.
 * @param {string} [main.warning_background] - The warning background color.
 * @param {string} [main.success] - The success color.
 * @param {string} [main.success_background] - The success background color.
 * @param {string} [main.info] - The info color.
 * @param {string} [main.info_background] - The info background color.
 * @param {number} [amount] - The amount to adjust the background color if needed.
 * @returns {Object} - An object containing standardized colors and their backgrounds.
 */
export const convertTypeColors = (main, amount) => ({
  error: main.error || COLORS.error,
  error_background: !main.error && !main.error_background
    ? COLORS.errorBackground
    : main.error_background || getLightenOrDarkenColor(main.error, amount),

  warning: main.warning || COLORS.warning,
  warning_background: !main.warning && !main.warning_background
    ? COLORS.warningBackground
    : main.warning_background || getLightenOrDarkenColor(main.warning, amount),

  success: main.success || COLORS.success,
  success_background: !main.success && !main.success_background
    ? COLORS.successBackground
    : main.success_background || getLightenOrDarkenColor(main.success, amount),

  info: main.info || COLORS.info,
  info_background: !main.info && !main.info_background
    ? COLORS.infoBackground
    : main.info_background || getLightenOrDarkenColor(main.info, amount),
});

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
  background: main.background ?? COLORS.background,
  active_color: main.active_color ?? COLORS.activeColor,

  ...convertTypeColors(main),
});

/**
 * Convert theme main colors to form object
 *
 * @param {ThemeTableColors} [table = {}]
 * @returns {ThemeTableColorsForm}
 */
export const themeTableColorsToForm = (table = {}) => {
  const isDarkRowColor = table.row_color && isDarkColor(table.row_color);
  const shiftRowColor = table.shift_row_color ?? (
    isDarkRowColor ? COLORS.table.shiftRowDarkColor : COLORS.table.shiftRowColor
  );
  const hoverRowColor = table.hover_row_color ?? (
    isDarkRowColor ? COLORS.table.hoverRowDarkColor : COLORS.table.hoverRowColor
  );

  return {
    background: table.background ?? COLORS.table.background,
    row_color: table.row_color ?? COLORS.table.rowColor,
    shift_row_color: {
      enabled: !!table.shift_row_color,
      color: shiftRowColor,
    },
    hover_row_color: {
      enabled: !!table.hover_row_color,
      color: hoverRowColor,
    },
  };
};

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
