import { getSecondByUnit } from '@/helpers/getSecondByUnit';

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

        const interval = parseInt(periodicRefresh.interval, 10);
        const oldInterval = parseInt(oldPeriodicRefresh.interval, 10);

        if (periodicRefresh.enabled && periodicRefresh.interval) {
          if (interval !== oldInterval || periodicRefresh.enabled !== oldPeriodicRefresh.enabled) {
            const delay = getSecondByUnit(interval, periodicRefresh.unit);

            if (this.periodicRefreshInterval) {
              clearInterval(this.periodicRefreshInterval);
            }

            this.periodicRefreshInterval = setInterval(() => this.fetchList({ isPeriodicRefresh: true }), delay * 1000);
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
