import { DEFAULT_TIMEZONE, DEFAULT_USER_INTERFACE_DATE_FORMAT } from '@/constants';

import { registerUserInterfaceDateFormatProvider } from '@/helpers/date/date-format';

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
        dateFormat: this.dateFormat || DEFAULT_USER_INTERFACE_DATE_FORMAT,
        dark: false,
        theme: {},
        setTheme: this.setTheme,
      },
    };
  },

  watch: {
    timezone: 'syncSystemInfo',
    dateFormat: 'syncSystemInfo',
  },

  created() {
    registerUserInterfaceDateFormatProvider(() => this.system.dateFormat);
  },

  methods: {
    syncSystemInfo() {
      this.setSystemData({
        timezone: this.timezone ?? DEFAULT_TIMEZONE,
        dateFormat: this.dateFormat || DEFAULT_USER_INTERFACE_DATE_FORMAT,
      });
    },

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
