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
import { themePropertiesToCSSVariables } from '@/helpers/entities/theme/entity';
import {
  convertTypeColors,
  tableColorsToTheme,
  themeToForm,
} from '@/helpers/entities/theme/form';

const DEFAULT_TEXT_ON_LIGHT = '#000000';
const DEFAULT_TEXT_ON_DARK = '#FFFFFF';

/**
 * Normalize a stored theme or a flat color object into a theme form.
 *
 * @param {Theme|Object} [theme={}] - Theme entity, theme form, or flat color map.
 * @returns {ThemeForm}
 */
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
      widget: flatColors.widget,
      state: flatColors.state,
      icons: flatColors.icons,
      status: flatColors.status,
    },
  });
};

/**
 * Build a readable text color for each state color.
 *
 * @param {ThemeStatesColors} state - State colors keyed by state name.
 * @returns {Object<string, string>} Text colors keyed as `{state}-text`.
 */
const createStateTextColors = state =>
  Object.entries(state).reduce((acc, [key, color]) => {
    acc[`${key}-text`] = getMostReadableTextColor(color);

    return acc;
  }, {});

const COLOR_VARIATION_INDEXES = [1, 2, 3, 4, 5];
const VUETIFY_PALETTE_VARIATION_KEYS = ['primary', 'secondary'];

/**
 * Build Vuetify darken and lighten palette steps for a hex color.
 *
 * @param {string} key - Palette color name.
 * @param {string} value - Source color. Non-hex values produce no variations.
 * @returns {Object<string, string>} Colors keyed as `{key}-darken-{n}` and `{key}-lighten-{n}`.
 */
const createPaletteColorVariations = (key, value) => {
  if (!String(value).startsWith('#')) {
    return {};
  }

  return COLOR_VARIATION_INDEXES.reduce((acc, index) => {
    acc[`${key}-darken-${index}`] = colorToHex(getDarkenColor(value, index * 10));
    acc[`${key}-lighten-${index}`] = colorToHex(getLightenColor(value, index * 10));

    return acc;
  }, {});
};

/**
 * Build the Vuetify theme color map from CSS variables and main palette colors.
 *
 * @param {Object} colors
 * @param {Object<string, string>} colors.variables - Theme CSS variables.
 * @param {ThemeMainColors} colors.main - Main palette colors.
 * @returns {Object<string, string>} Hex colors, including primary and secondary variations.
 */
const createVuetifyColors = ({ variables, main }) => {
  const vuetifyVariables = Object.fromEntries(
    Object.entries(variables).filter(([key]) => !key.startsWith('table-')),
  );
  const colors = Object.entries({
    ...vuetifyVariables,
    background: variables['application-background'] ?? main.background,
    surface: main.background,
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

  return VUETIFY_PALETTE_VARIATION_KEYS.reduce(
    (acc, key) => ({
      ...acc,
      ...createPaletteColorVariations(key, acc[key]),
    }),
    colors,
  );
};

/**
 * Build text, button, stepper, and divider variables for the active theme.
 *
 * @param {Object} options
 * @param {string} options.activeColor - Color used for active controls.
 * @param {string} options.darkBaseColor - Base color for text on dark surfaces.
 * @param {boolean} options.isDark - Whether the theme background is dark.
 * @param {string} options.lightBaseColor - Base color for text on light surfaces.
 * @param {Object<string, string>} options.typeBackgroundColors - Error, warning, success, and info backgrounds.
 * @returns {Object} Runtime theme variables.
 */
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

/**
 * Build legacy Vuetify 2 darken and lighten CSS declarations for a hex color.
 *
 * @param {string} key - CSS variable name without the `--v-` prefix.
 * @param {string} value - Source color. Non-hex values produce no variations.
 * @returns {string[]} Declarations for hyphenated and non-hyphenated variation names.
 */
const createLegacyColorVariations = (key, value) => {
  if (!String(value).startsWith('#')) {
    return [];
  }

  return [1, 2, 3, 4, 5].flatMap(index => [
    `--v-${key}-darken${index}:${getDarkenColor(value, index * 10)};`,
    `--v-${key}-darken-${index}:${getDarkenColor(value, index * 10)};`,
    `--v-${key}-lighten${index}:${getLightenColor(value, index * 10)};`,
    `--v-${key}-lighten-${index}:${getLightenColor(value, index * 10)};`,
  ]);
};

/**
 * Join legacy `--v-*` and table `--c-*` CSS variables into one declaration string.
 *
 * @param {Object} variables
 * @param {Object<string, string>} variables.themeVariables - Theme CSS variables.
 * @param {Object} variables.runtimeVariables - Runtime text, button, and divider variables.
 * @returns {string}
 */
const createLegacyCssVariables = ({ themeVariables, runtimeVariables }) => {
  const legacyThemeVariables = Object.entries(themeVariables).reduce((acc, [key, value]) => {
    if (key.startsWith('table-')) {
      return acc;
    }

    acc.push(`--v-${key}:${value};`);
    acc.push(`--v-${key}-base:${value};`);
    acc.push(...createLegacyColorVariations(key, value));

    return acc;
  }, []);

  const legacyRuntimeVariables = Object.entries(
    themePropertiesToCSSVariables(runtimeVariables),
  ).map(([key, value]) => `--v-${key}:${value};`);

  const tableCssVariables = Object.entries(themeVariables)
    .filter(([key]) => key.startsWith('table-'))
    .map(([key, value]) => `--c-${key}:${value};`);

  return [...legacyThemeVariables, ...legacyRuntimeVariables, ...tableCssVariables].join('');
};

/**
 * Convert a Canopsis theme into a Vuetify theme and the legacy CSS still used by the app.
 *
 * @param {Theme|Object} [theme={}] - Theme entity, theme form, or flat color map.
 * @returns {{
 *   definition: Object,
 *   isDark: boolean,
 *   legacyCssVariables: string,
 *   runtimeVariables: Object,
 *   theme: ThemeForm,
 * }}
 */
export const createCanopsisVuetifyTheme = (theme = {}) => {
  const form = getThemeForm(theme);
  const {
    main,
    table,
    widget,
    state,
    icons,
    status,
  } = form.colors;
  const isDark = isDarkColor(main.background);
  const stateText = createStateTextColors(state);
  const typeColors = convertTypeColors(main);
  const tableVariables = tableColorsToTheme(table);
  const vuetifyVariables = merge({}, DEFAULT_THEME_COLORS, {
    table: tableVariables,
    widget: {
      title_background: widget.title_background,
      content_background: widget.content_background,
      border: widget.border.enabled,
      border_color: widget.border.color,
      border_radius: widget.border.radius,
    },
    'impact-state': COLORS.impactStateGrouped,
    state: {
      ...state,
      ...stateText,
    },
    icons,
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
