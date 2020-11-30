import Vue from 'vue';

import { DEFAULT_TIMEZONE } from '@/constants';
import { DEFAULT_JOB_EXECUTOR_FETCH_TIMEOUT_SECONDS } from '@/config';

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
        jobExecutorFetchTimeoutSeconds: this.jobExecutorFetchTimeoutSeconds
          || DEFAULT_JOB_EXECUTOR_FETCH_TIMEOUT_SECONDS,
      },
    };
  },
  methods: {
    /**
     * @param {Object} options
     * @param {string} [options.timezone]
     * @param {number} [options.jobExecutorFetchTimeoutSeconds]
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
