import Vue from 'vue';

import { DEFAULT_TIMEZONE } from '@/constants';
import { INSTRUCTION_EXECUTE_JOB_ALERT_DELAY } from '@/config';

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
        jobExecutorFetchTimeout: this.jobExecutorFetchTimeout || INSTRUCTION_EXECUTE_JOB_ALERT_DELAY,
      },
    };
  },
  methods: {
    /**
     * @param {Object} options
     * @param {string} [options.timezone]
     * @param {number} [options.jobExecutorFetchTimeout]
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
