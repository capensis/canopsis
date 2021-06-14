export const widgetPeriodicRefreshMixin = {
  data() {
    return {
      periodicRefreshInterval: null,
    };
  },
  watch: {
    'widget.parameters.periodic_refresh': {
      immediate: true,
      handler(value, oldValue) {
        const periodicRefresh = value || {};
        const oldPeriodicRefresh = oldValue || {};

        if (periodicRefresh.enabled && periodicRefresh.seconds) {
          const secondsIsChanged = periodicRefresh.seconds !== oldPeriodicRefresh.seconds;
          const enabledIsChanged = periodicRefresh.enabled !== oldPeriodicRefresh.enabled;

          if (secondsIsChanged || enabledIsChanged) {
            if (this.periodicRefreshInterval) {
              clearInterval(this.periodicRefreshInterval);
            }

            this.periodicRefreshInterval = setInterval(() => {
              this.fetchList({ isPeriodicRefresh: true });
            }, periodicRefresh.seconds * 1000);
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
