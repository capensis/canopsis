import Vue from 'vue';
import theme from 'vuetify/es5/components/Vuetify/mixins/theme';

import { DEFAULT_TIMEZONE } from '@/constants';

import { themeColorsToCSSVariables } from '@/helpers/entities/theme/entity';
import { getDarkenColor, isDarkColor } from '@/helpers/color';

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

    setTheme({ colors }) {
      const { main, table, state } = colors;
      const variables = themeColorsToCSSVariables({
        ...main,
        table,
        state,
        applicationBackground: getDarkenColor(main.background, 5),
      });

      this.$vuetify.theme = theme(variables);

      this.system.dark = isDarkColor(main.background);
    },
  },
};
