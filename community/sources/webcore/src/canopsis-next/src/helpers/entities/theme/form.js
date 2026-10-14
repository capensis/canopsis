import { COLORS } from '@/config';
import { DEFAULT_WIDGET_BORDER_RADIUS, THEME_FONT_SIZES } from '@/constants';

import { getLightenOrDarkenColor, isDarkColor } from '@/helpers/color';

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
 * @property {ThemeEnabledColor} shift_row
 * @property {ThemeEnabledColor} hover_row
 */

/**
 * @typedef {Object} ThemeWidgetColorsForm
 * @property {string} title_background
 * @property {string} content_background
 * @property {ThemeEnabledColor} border
 */

/**
 * @typedef {ThemeStatesColors} ThemeStatesColorsForm
 */

/**
 * @typedef {ThemeIconsColors} ThemeIconsColorsForm
 */

/**
 * @typedef {ThemeColors} ThemeColorsForm
 * @property {ThemeMainColorsForm} main
 * @property {ThemeTableColorsForm} table
 * @property {ThemeWidgetColorsForm} widget
 * @property {ThemeStatesColorsForm} state
 * @property {ThemeIconsColorsForm} icons
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
export const convertTypeColors = (main = {}, amount) => ({
  error: main.error || COLORS.error,
  error_background:
    !main.error && !main.error_background
      ? COLORS.errorBackground
      : main.error_background || getLightenOrDarkenColor(main.error, amount),

  warning: main.warning || COLORS.warning,
  warning_background:
    !main.warning && !main.warning_background
      ? COLORS.warningBackground
      : main.warning_background || getLightenOrDarkenColor(main.warning, amount),

  success: main.success || COLORS.success,
  success_background:
    !main.success && !main.success_background
      ? COLORS.successBackground
      : main.success_background || getLightenOrDarkenColor(main.success, amount),

  info: main.info || COLORS.info,
  info_background:
    !main.info && !main.info_background
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
  primary: main.primary || COLORS.primary,
  secondary: main.secondary || COLORS.secondary,
  accent: main.accent || COLORS.accent,
  background: main.background || COLORS.background,
  active_color: main.active_color || COLORS.activeColor,

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
  const shiftRowEnabled = table.shift_row?.enabled ?? Boolean(table.shift_row);
  const hoverRowEnabled = table.hover_row?.enabled ?? Boolean(table.hover_row);
  const shiftRowColor =
    table.shift_row_color ||
    table.shift_row?.color ||
    (isDarkRowColor ? COLORS.table.shiftRowDarkColor : COLORS.table.shiftRowColor);
  const hoverRowColor =
    table.hover_row_color ||
    table.hover_row?.color ||
    (isDarkRowColor ? COLORS.table.hoverRowDarkColor : COLORS.table.hoverRowColor);

  return {
    background: table.background || COLORS.table.background,
    row_color: table.row_color || COLORS.table.rowColor,
    shift_row: {
      enabled: shiftRowEnabled,
      color: shiftRowColor,
    },
    hover_row: {
      enabled: hoverRowEnabled,
      color: hoverRowColor,
    },
  };
};

/**
 * Convert theme widget colors to form object.
 *
 * @param {ThemeWidgetColors} [widget = {}]
 * @returns {ThemeWidgetColorsForm}
 */
export const themeWidgetColorsToForm = (widget = {}) => ({
  border: {
    enabled: widget.border?.enabled ?? Boolean(widget.border),
    color: widget.border_color || widget.border?.color || COLORS.widget.borderColor,
    radius: widget.border_radius ?? widget.border?.radius ?? DEFAULT_WIDGET_BORDER_RADIUS,
  },
  title_background: widget.title_background || COLORS.widget.titleBackground,
  content_background: widget.content_background || COLORS.widget.contentBackground,
});

/**
 * Convert theme main colors to form object
 *
 * @param {ThemeStatesColors} [states = {}]
 * @returns {ThemeStatesColorsForm}
 */
export const themeStatesColorsToForm = (states = {}) => ({
  ok: states.ok || COLORS.state.ok,
  minor: states.minor || COLORS.state.minor,
  major: states.major || COLORS.state.major,
  critical: states.critical || COLORS.state.critical,
  pause: states.pause || COLORS.state.pause,
});

/**
 * Convert extra details colors to form object
 *
 * @param {ThemeIconsColors} [icons = {}]
 * @returns {ThemeIconsColorsForm}
 */
export const themeIconsColorsToForm = (icons = {}) => ({
  ack: icons.ack || COLORS.alarmExtraDetails.ack,
  canceled: icons.canceled || COLORS.alarmExtraDetails.canceled,
  children: icons.children || COLORS.alarmExtraDetails.children,
  parent: icons.parent || COLORS.alarmExtraDetails.parent,
  comment: icons.comment || COLORS.alarmExtraDetails.comment,
  snooze: icons.snooze || COLORS.alarmExtraDetails.snooze,
  tickets: icons.tickets || icons.ticket || COLORS.alarmExtraDetails.tickets,
  cancel_delay:
    icons.cancel_delay ||
    icons.close_delay ||
    icons.closeDelay ||
    COLORS.alarmExtraDetails.cancel_delay,
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
  widget: themeWidgetColorsToForm(colors.widget),
  state: themeStatesColorsToForm(colors.state),
  icons: themeIconsColorsToForm(colors.icons ?? colors.extra_details),
  status: colors.status || COLORS.status,
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
 * Convert table colors to theme object.
 *
 * @param {ThemeTableColors|ThemeTableColorsForm} [table = {}]
 * @returns {ThemeTableColors}
 */
export const tableColorsToTheme = (table = {}) => {
  const shiftRowEnabled = table.shift_row?.enabled ?? Boolean(table.shift_row);
  const hoverRowEnabled = table.hover_row?.enabled ?? Boolean(table.hover_row);

  return {
    ...table,
    shift_row: shiftRowEnabled,
    shift_row_color: shiftRowEnabled
      ? table.shift_row?.color || table.shift_row_color || null
      : null,
    hover_row: hoverRowEnabled,
    hover_row_color: hoverRowEnabled
      ? table.hover_row?.color || table.hover_row_color || null
      : null,
  };
};

/**
 * Convert widget colors from form to API object.
 *
 * @param {ThemeWidgetColorsForm} [widget = {}]
 * @returns {ThemeWidgetColors}
 */
export const formWidgetColorsToTheme = (widget = {}) => {
  const preparedWidget = themeWidgetColorsToForm(widget);

  return {
    title_background: preparedWidget.title_background,
    content_background: preparedWidget.content_background,
    border: preparedWidget.border.enabled,
    border_color: preparedWidget.border.color,
    border_radius: preparedWidget.border.radius,
  };
};

/**
 * Convert colors from form to API object.
 *
 * @param {ThemeColorsForm} [colors = {}]
 * @returns {ThemeColors}
 */
export const formColorsToTheme = (colors = {}) => {
  const { extra_details: extraDetails, ...preparedColors } = colors;

  return {
    ...preparedColors,
    table: tableColorsToTheme(colors.table),
    widget: formWidgetColorsToTheme(colors.widget),
    icons: themeIconsColorsToForm(colors.icons ?? extraDetails),
  };
};

/**
 * Convert theme form to API object.
 *
 * @param {ThemeForm} form
 * @returns {Theme}
 */
export const formToTheme = form => ({
  ...form,
  colors: formColorsToTheme(form.colors),
});
