import { merge } from 'lodash';

import { COLORS, DEFAULT_THEME_COLORS } from '@/config';
import { THEME_FONT_PIXEL_SIZES, THEME_FONT_SIZES } from '@/constants';

import {
  colorToRgba,
  colorToHex,
  getDarkenColor,
  getLightenColor,
  getMostReadableTextColor,
  isDarkColor,
} from '@/helpers/color';
import { convertTypeColors, themeToForm } from '@/helpers/entities/theme/form';
import { themePropertiesToCSSVariables } from '@/helpers/entities/theme/entity';

const DEFAULT_TEXT_ON_LIGHT = '#000000';
const DEFAULT_TEXT_ON_DARK = '#FFFFFF';

const getThemeForm = (theme = {}) => {
  if (theme.colors?.main) {
    return themeToForm(theme);
  }

  const flatColors = theme.colors ?? theme;

  return themeToForm({
    name: theme.name,
    font_size: theme.font_size,
    colors: {
      main: flatColors,
      table: flatColors.table,
      state: flatColors.state,
      status: flatColors.status,
    },
  });
};

const createStateTextColors = (state) =>
  Object.entries(state).reduce((acc, [key, color]) => {
    acc[`${key}-text`] = getMostReadableTextColor(color);

    return acc;
  }, {});

const createVuetifyColors = ({ variables, main, table }) =>
  Object.entries({
    ...variables,
    background: variables['application-background'] ?? main.background,
    surface: table.background ?? main.background,
    primary: main.primary,
    secondary: main.secondary,
    accent: main.accent,
    error: main.error,
    warning: main.warning,
    success: main.success,
    info: main.info,
  }).reduce((acc, [key, value]) => {
    acc[key] = colorToHex(value);

    return acc;
  }, {});

const createVuetifyRuntimeVariables = ({
  activeColor,
  darkBaseColor,
  isDark,
  lightBaseColor,
  typeBackgroundColors,
}) => {
  const textLight = {
    primary: colorToRgba(lightBaseColor, 0.87),
    secondary: colorToRgba(lightBaseColor, 0.54),
    disabled: colorToRgba(lightBaseColor, 0.38),
  };
  const textDark = {
    primary: darkBaseColor,
    secondary: colorToRgba(darkBaseColor, 0.7),
    disabled: colorToRgba(darkBaseColor, 0.5),
  };
  const buttonsLight = {
    disabled: colorToRgba(lightBaseColor, 0.26),
    focused: colorToRgba(lightBaseColor, 0.12),
  };
  const buttonsDark = {
    disabled: colorToRgba(lightBaseColor, 0.3),
    focused: colorToRgba(lightBaseColor, 0.12),
  };
  const active = isDark ? DEFAULT_TEXT_ON_LIGHT : DEFAULT_TEXT_ON_DARK;
  const completed = isDark ? DEFAULT_TEXT_ON_DARK : DEFAULT_TEXT_ON_LIGHT;
  const stepper = {
    active,
    completed: colorToRgba(completed, 0.87),
    hover: colorToRgba(completed, isDark ? 0.75 : 0.54),
  };
  const divider = {
    borderColor: colorToRgba(isDark ? darkBaseColor : lightBaseColor, 0.12),
  };

  return {
    textLight,
    textDark,
    buttonsLight,
    buttonsDark,
    stepper,
    divider,
    activeColor,

    ...typeBackgroundColors,
  };
};

const createLegacyColorVariations = (key, value) => {
  if (!String(value).startsWith('#')) {
    return [];
  }

  return [1, 2, 3, 4, 5].flatMap((index) => [
    `--v-${key}-darken${index}:${getDarkenColor(value, index * 5)};`,
    `--v-${key}-darken-${index}:${getDarkenColor(value, index * 5)};`,
    `--v-${key}-lighten${index}:${getLightenColor(value, index * 5)};`,
    `--v-${key}-lighten-${index}:${getLightenColor(value, index * 5)};`,
  ]);
};

const createLegacyCssVariables = ({ themeVariables, runtimeVariables }) => {
  const legacyThemeVariables = Object.entries(themeVariables).reduce((acc, [key, value]) => {
    acc.push(`--v-${key}:${value};`);
    acc.push(`--v-${key}-base:${value};`);
    acc.push(...createLegacyColorVariations(key, value));

    return acc;
  }, []);

  const legacyRuntimeVariables = Object.entries(
    themePropertiesToCSSVariables(runtimeVariables),
  ).map(([key, value]) => `--v-${key}:${value};`);

  return [...legacyThemeVariables, ...legacyRuntimeVariables].join('');
};

export const createCanopsisVuetifyTheme = (theme = {}) => {
  const form = getThemeForm(theme);
  const { main, table, state, status } = form.colors;
  const isDark = isDarkColor(main.background);
  const stateText = createStateTextColors(state);
  const typeColors = convertTypeColors(main);
  const vuetifyVariables = merge({}, DEFAULT_THEME_COLORS, {
    table,
    'impact-state': COLORS.impactStateGrouped,
    state: {
      ...state,
      ...stateText,
    },
    status,
    applicationBackground: getDarkenColor(main.background, isDark ? 7 : 2),

    ...main,
    ...typeColors,
  });

  const themeVariables = themePropertiesToCSSVariables(vuetifyVariables);
  const lightBaseColor = isDark ? DEFAULT_TEXT_ON_LIGHT : main.active_color;
  const darkBaseColor = isDark ? main.active_color : DEFAULT_TEXT_ON_DARK;
  const typeBackgroundColors = {
    'error-background': themeVariables['error-background'],
    'warning-background': themeVariables['warning-background'],
    'success-background': themeVariables['success-background'],
    'info-background': themeVariables['info-background'],
  };
  const runtimeVariables = {
    fontSizeRoot: `${THEME_FONT_PIXEL_SIZES[form.font_size] ?? THEME_FONT_PIXEL_SIZES[THEME_FONT_SIZES.medium]}px`,

    ...createVuetifyRuntimeVariables({
      activeColor: main.active_color,
      darkBaseColor,
      isDark,
      lightBaseColor,
      typeBackgroundColors,
    }),
  };
  const definition = {
    dark: isDark,
    colors: createVuetifyColors({
      variables: themeVariables,
      main: {
        ...main,
        ...typeColors,
      },
      table,
    }),
    variables: {
      'border-color': isDark ? DEFAULT_TEXT_ON_DARK : DEFAULT_TEXT_ON_LIGHT,
      'border-opacity': 0.12,
      'high-emphasis-opacity': isDark ? 1 : 0.87,
      'medium-emphasis-opacity': isDark ? 0.7 : 0.6,
      'disabled-opacity': isDark ? 0.5 : 0.38,
      'theme-on-dark': DEFAULT_TEXT_ON_DARK,
      'theme-on-light': DEFAULT_TEXT_ON_LIGHT,
    },
  };

  return {
    definition,
    isDark,
    legacyCssVariables: createLegacyCssVariables({
      themeVariables,
      runtimeVariables,
    }),
    runtimeVariables,
    theme: form,
  };
};
