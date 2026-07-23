import { DEFAULT_THEME_COLORS } from '@/config';
import { createCanopsisVuetifyTheme } from '@/helpers/entities/theme/vuetify';
import { CANOPSIS_THEME_NAME, vuetify } from '@/plugins/vuetify';

const DEFAULT_CANOPSIS_THEME = createCanopsisVuetifyTheme(DEFAULT_THEME_COLORS);

export const systemThemeMixin = {
  data() {
    return {
      themeCssVariables: DEFAULT_CANOPSIS_THEME.legacyCssVariables,
    };
  },

  computed: {
    css() {
      return `:root{${this.themeCssVariables}}`;
    },
  },

  created() {
    this.injectCSS(this.css);
  },

  beforeUnmount() {
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
        this.styleElement.id = 'canopsis-theme-legacy-variables';
      }

      this.updateCSS(value);
      document.head.appendChild(this.styleElement);
    },

    updateCSS(value) {
      if (this.styleElement) {
        this.styleElement.innerHTML = value;
      }
    },

    ejectCSS() {
      this.styleElement?.remove();
      delete this.styleElement;
    },

    setTheme(theme = DEFAULT_THEME_COLORS) {
      const {
        definition,
        isDark,
        legacyCssVariables,
        theme: preparedTheme,
      } = createCanopsisVuetifyTheme(theme);

      vuetify.theme.themes.value = {
        ...vuetify.theme.themes.value,
        [CANOPSIS_THEME_NAME]: definition,
      };

      vuetify.theme.change(CANOPSIS_THEME_NAME);

      this.themeCssVariables = legacyCssVariables;
      this.system.dark = isDark;
      this.system.theme = preparedTheme;
    },
  },
};

export default systemThemeMixin;
