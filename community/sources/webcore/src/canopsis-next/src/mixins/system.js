import Vue from 'vue';

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
  },
};
