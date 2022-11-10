import Vue from 'vue';

import theme from 'vuetify/es5/components/Vuetify/mixins/theme';
import { THEMES } from '@/config';

import { DEFAULT_TIMEZONE } from '@/constants';

export const systemMixin = {
  provide() {
    return {
      $system: this.system,
    };
  },
  data() {
    return {
      system: {
        timezone: this.timezone || DEFAULT_TIMEZONE,
        darkMode: false,
        themes: Object.keys(THEMES),
        setTheme: this.setTheme,
        setDarkMode: this.setDarkMode,
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

    setTheme(name) {
      const { dark, colors } = THEMES[name];

      this.$vuetify.theme = theme(colors);

      this.system.darkMode = dark;
    },
  },
};
