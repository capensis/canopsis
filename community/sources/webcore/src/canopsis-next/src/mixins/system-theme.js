import { kebabCase, merge } from 'lodash';

import { DEFAULT_THEME_COLORS } from '@/config';

import { colorToRgba, getDarkenColor, isDarkColor } from '@/helpers/color';
import { themePropertiesToCSSVariables } from '@/helpers/entities/theme/entity';

import { THEME_FONT_PIXEL_SIZES, THEME_FONT_SIZES } from '@/constants/theme';

export const systemThemeMixin = {
  data() {
    return {
      otherVariables: {
        fontSizeRoot: `${THEME_FONT_PIXEL_SIZES[THEME_FONT_SIZES.medium]}px`,
      },
    };
  },
  computed: {
    cssVariables() {
      return Object.entries(themePropertiesToCSSVariables(this.otherVariables))
        .reduce((acc, [key, value]) => `${acc}--v-${kebabCase(key)}:${value};`, '');
    },
    css() {
      return `:root{${this.cssVariables}}`;
    },
  },
  created() {
    this.injectCSS(this.css);
  },
  beforeDestroy() {
    this.ejectCSS();
  },
  watch: {
    css(value) {
      this.updateCSS(value);
    },
  },
  methods: {
    injectCSS(value) {
      if (!this.styleElement) {
        this.styleElement = document.createElement('style');
      }

      this.updateCSS(value);

      document.head.appendChild(this.styleElement);
    },

    updateCSS(value) {
      this.styleElement.innerHTML = value;
    },

    ejectCSS() {
      document.head.removeChild(this.styleElement);
    },

    setTheme(theme) {
      const { colors, font_size: fontSize } = theme;
      const { main, table, state } = colors;

      const white = '#fff';
      const black = '#000';

      const isDark = isDarkColor(main.background);

      const vuetifyVariables = merge({}, DEFAULT_THEME_COLORS, {
        ...main,
        table,
        state,
        applicationBackground: getDarkenColor(main.background, isDark ? 7 : 2),
      });

      const variables = themePropertiesToCSSVariables(vuetifyVariables);

      this.$vuetify.theme.dark = isDark;
      this.$vuetify.theme.themes.dark = variables;
      this.$vuetify.theme.themes.light = variables;

      const lightBaseColor = isDark ? black : main.active_color;
      const darkBaseColor = isDark ? main.active_color : white;

      const textLight = {
        primary: colorToRgba(lightBaseColor, 0.87),
        secondary: colorToRgba(lightBaseColor, 0.54),
        disabled: colorToRgba(lightBaseColor, 0.38),
      };
      const textDark = {
        primary: darkBaseColor,
        secondary: colorToRgba(darkBaseColor, 0.70),
        disabled: colorToRgba(darkBaseColor, 0.50),
      };

      const buttonsLight = {
        disabled: colorToRgba(lightBaseColor, 0.26),
        focused: colorToRgba(lightBaseColor, 0.12),
      };
      const buttonsDark = {
        disabled: colorToRgba(lightBaseColor, 0.3),
        focused: colorToRgba(lightBaseColor, 0.12),
      };

      const active = isDark ? black : white;
      const completed = isDark ? white : black;

      const stepper = {
        active,
        completed: colorToRgba(completed, 0.87),
        hover: colorToRgba(completed, isDark ? 0.75 : 0.54),
      };

      this.otherVariables = {
        fontSizeRoot: `${THEME_FONT_PIXEL_SIZES[fontSize]}px`,
        textLight,
        textDark,
        buttonsLight,
        buttonsDark,
        stepper,
      };
      this.system.dark = isDark;
      this.system.theme = theme;
    },
  },
};
