import Vue from 'vue';

import { DEFAULT_TIMEZONE } from '@/constants';

export default {
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
    setSystemData(options) {
      Object.entries(options).forEach(([key, value]) => {
        Vue.set(this.system, key, value);
      });
    },
  },
};
