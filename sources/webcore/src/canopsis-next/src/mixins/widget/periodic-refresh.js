export default {
  data() {
    return {
      periodicRefreshInterval: null,
    };
  },
  watch: {
    'widget.parameters.periodicRefresh': {
      immediate: true,
      handler(value, oldValue) {
        const periodicRefresh = value || {};
        const oldPeriodicRefresh = oldValue || {};

        const interval = parseInt(periodicRefresh.interval, 10) * 1000; // In seconds
        const oldInterval = parseInt(oldPeriodicRefresh.interval, 10) * 1000; // In seconds

        if (periodicRefresh.enabled && periodicRefresh.interval) {
          if (interval !== oldInterval || periodicRefresh.enabled !== oldPeriodicRefresh.enabled) {
            if (this.periodicRefreshInterval) {
              clearInterval(this.periodicRefreshInterval);
            }

            this.periodicRefreshInterval = setInterval(() => this.fetchList({ isPeriodicRefresh: true }), interval);
          }
        } else {
          clearInterval(this.periodicRefreshInterval);
        }
      },
    },
  },
  beforeDestroy() {
    if (this.periodicRefreshInterval) {
      clearInterval(this.periodicRefreshInterval);
    }
  },
};
