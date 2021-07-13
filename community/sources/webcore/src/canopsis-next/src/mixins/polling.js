/**
 * Mixin creator for polling
 *
 * @param {string} method
 * @param {boolean} [startOnMount = false]
 * @param {number} delay
 * @returns {{
 *   data(): { timeout: null },
 *   methods: {
 *     startPolling(): void,
 *     stopPolling(): void
 *   },
 *   beforeDestroy(): void,
 *   mounted(): void
 * }}
 */
export const createPollingMixin = ({ method, startOnMount = false, delay }) => ({
  data() {
    return {
      timeout: null,
    };
  },
  mounted() {
    if (startOnMount) {
      this.startPolling();
    }
  },
  beforeDestroy() {
    this.stopPolling();
  },
  methods: {
    async polling() {
      await this[method]();

      this.startPolling();
    },

    startPolling() {
      this.timeout = setTimeout(this.polling, delay);
    },

    stopPolling() {
      clearTimeout(this.timeout);

      this.timeout = null;
    },
  },
});
