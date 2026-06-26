import { COLORS } from '@/config'
import { THEME_FONT_SIZES } from '@/constants'

import { getLightenOrDarkenColor, isDarkColor } from '@/helpers/color'

export const convertTypeColors = (main = {}, amount) => ({
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
})

export const themeMainColorsToForm = (main = {}) => ({
  primary: main.primary ?? COLORS.primary,
  secondary: main.secondary ?? COLORS.secondary,
  accent: main.accent ?? COLORS.accent,
  background: main.background ?? COLORS.background,
  active_color: main.active_color ?? COLORS.activeColor,

  ...convertTypeColors(main),
})

export const themeTableColorsToForm = (table = {}) => {
  const isDarkRowColor = table.row_color && isDarkColor(table.row_color)
  const shiftRowColor = table.shift_row_color ?? (
    isDarkRowColor ? COLORS.table.shiftRowDarkColor : COLORS.table.shiftRowColor
  )
  const hoverRowColor = table.hover_row_color ?? (
    isDarkRowColor ? COLORS.table.hoverRowDarkColor : COLORS.table.hoverRowColor
  )

  return {
    background: table.background ?? COLORS.table.background,
    row_color: table.row_color ?? COLORS.table.rowColor,
    shift_row: {
      enabled: Boolean(table.shift_row),
      color: shiftRowColor,
    },
    hover_row: {
      enabled: Boolean(table.hover_row),
      color: hoverRowColor,
    },
  }
}

export const themeStatesColorsToForm = (states = {}) => ({
  ok: states.ok ?? COLORS.state.ok,
  minor: states.minor ?? COLORS.state.minor,
  major: states.major ?? COLORS.state.major,
  critical: states.critical ?? COLORS.state.critical,
  pause: states.pause ?? COLORS.state.pause,
})

export const themeColorsToForm = (colors = {}) => ({
  main: themeMainColorsToForm(colors.main),
  table: themeTableColorsToForm(colors.table),
  state: themeStatesColorsToForm(colors.state),
  status: colors.status ?? COLORS.status,
})

export const themeToForm = (theme = {}) => ({
  name: theme.name ?? '',
  font_size: theme.font_size ?? THEME_FONT_SIZES.medium,
  colors: themeColorsToForm(theme.colors),
})
