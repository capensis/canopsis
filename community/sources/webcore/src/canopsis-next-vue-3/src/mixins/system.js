import { DEFAULT_TIMEZONE } from '@/constants';

import { systemThemeMixin } from '@/mixins/system-theme';

export const systemMixin = {
  mixins: [systemThemeMixin],

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
        theme: {},
        setTheme: this.setTheme,
      },
    };
  },

  methods: {
    setSystemData(options = {}) {
      Object.entries(options).forEach(([key, value]) => {
        if (value !== undefined) {
          this.system[key] = value;
        }
      });
    },
  },
};

export default systemMixin;
