import theme from 'vuetify/es5/components/Vuetify/mixins/theme';

import { getDarkenColor, isDarkColor } from '@/helpers/color';
import { themeColorsToCSSVariables } from '@/helpers/entities/theme/entity';

import { THEME_FONT_PIXEL_SIZES, THEME_FONT_SIZES } from '@/constants/theme';

export const systemThemeMixin = {
  data() {
    return {
      fontSize: THEME_FONT_PIXEL_SIZES[THEME_FONT_SIZES.medium],
    };
  },
  computed: {
    cssVariables() {
      return `--v-font-size-root:${this.fontSize}px;}`;
    },
    css() {
      return `:root{${this.cssVariables}`;
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

    setTheme({ colors }) {
      const { main, table, state } = colors;

      const isDark = isDarkColor(main.background);

      const variables = themeColorsToCSSVariables({
        ...main,
        table,
        state,
        applicationBackground: getDarkenColor(main.background, 7),
      });

      this.$vuetify.theme = theme(variables);

      this.fontSize = THEME_FONT_PIXEL_SIZES[main.font_size];
      this.system.dark = isDark;
    },
  },
};
