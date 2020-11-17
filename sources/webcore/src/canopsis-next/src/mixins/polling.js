/**
 * Mixin creator for polling
 *
 * @param {string} method
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
export default ({ method, delay }) => ({
  data() {
    return {
      timeout: null,
    };
  },
  mounted() {
    this.startPolling();
  },
  beforeDestroy() {
    this.stopPolling();
  },
  methods: {
    async startPolling() {
      await this[method]();

      this.timeout = setTimeout(this.startPolling, delay);
    },

    stopPolling() {
      clearTimeout(this.timeout);

      this.timeout = null;
    },
  },
});
