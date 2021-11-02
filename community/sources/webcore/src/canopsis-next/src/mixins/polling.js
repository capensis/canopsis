import { isFunction } from 'lodash';

/**
 * Mixin creator for polling
 *
 * @param {string} method
 * @param {string} [delayField = 'pollingDelay']
 * @returns {{
 *   methods: {
 *     polling(): Promise,
 *     stopPolling(): void,
 *     stopPolling(): void,
 *   },
 *   mounted(): void,
 *   beforeDestroy(): void,
 * }}
 */
export const pollingMixinCreator = ({ method, delayField = 'pollingDelay' }) => ({
  mounted() {
    this.startPolling();
  },
  beforeDestroy() {
    this.stopPolling();
  },
  methods: {
    async polling() {
      if (!isFunction(this[method])) {
        throw new Error(`Method ${method} not found`);
      }

      await this[method]();

      this.startPolling();
    },

    startPolling() {
      const delay = this[delayField];

      if (!delay) {
        return;
      }

      this.timeout = setTimeout(this.polling, delay);
    },

    stopPolling() {
      clearTimeout(this.timeout);

      this.timeout = null;
    },
  },
});
