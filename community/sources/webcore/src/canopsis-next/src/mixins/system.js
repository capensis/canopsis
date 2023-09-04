import Vue from 'vue';
import theme from 'vuetify/es5/components/Vuetify/mixins/theme';

import { DEFAULT_THEME_COLORS, THEMES, THEMES_NAMES } from '@/config';
import { DEFAULT_TIMEZONE } from '@/constants';

import { themeColorsToCSSVariables } from '@/helpers/entities/theme/entity';

export const systemMixin = {
  provide() {
    return {
      $system: this.system,
    };
  },
  data() {
    return {
      system: {
        timezone: this.timezone ?? DEFAULT_TIMEZONE,
        dark: false,
        setTheme: this.setTheme,
      },
    };
  },
  methods: {
    /**
     * @param {Object} options
     * @param {string} [options.timezone]
     */
    setSystemData(options) {
      Object.entries(options).forEach(([key, value]) => {
        if (value !== undefined) {
          Vue.set(this.system, key, value);
        }
      });
    },

    setTheme(name = THEMES_NAMES.canopsis) {
      if (THEMES[name]) {
        const { dark } = THEMES[name];
        const { colors } = {
          colors: {
            ...DEFAULT_THEME_COLORS,
          },
        };

        const variables = themeColorsToCSSVariables(colors);

        this.$vuetify.theme = theme(variables);

        this.system.dark = dark;
      }
    },
  },
};
