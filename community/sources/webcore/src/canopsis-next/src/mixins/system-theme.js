import theme from 'vuetify/es5/components/Vuetify/mixins/theme';
import { kebabCase } from 'lodash';

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
      this.styleElement = document.createElement('style');

      this.updateCSS(value);

      document.head.appendChild(this.styleElement);
    },

    updateCSS(value) {
      this.styleElement.innerHTML = value;
    },

    ejectCSS() {
      document.head.removeChild(this.styleElement);
    },

    setTheme({ colors, font_size: fontSize }) {
      const { main, table, state } = colors;

      const isDark = isDarkColor(main.background);

      const variables = themePropertiesToCSSVariables({
        ...main,
        table,
        state,
        applicationBackground: getDarkenColor(main.background, isDark ? 7 : 2),
      });

      this.$vuetify.theme = theme(variables);

      const lightBaseColor = isDark ? '#000' : main.active_color;
      const darkBaseColor = isDark ? main.active_color : '#fff';

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

      const stepperBaseColor = isDark ? '#000' : '#fff';
      const stepper = {
        active: stepperBaseColor,
        completed: colorToRgba(stepperBaseColor, 0.87),
        hover: colorToRgba(stepperBaseColor, isDark ? 0.75 : 0.54),
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
    },
  },
};
