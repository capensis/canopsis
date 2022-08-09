import { toSeconds } from '@/helpers/date/duration';

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
        const periodicRefresh = value;
        const oldPeriodicRefresh = oldValue ?? {};

        if (periodicRefresh?.enabled && periodicRefresh?.value) {
          const valueIsChanged = periodicRefresh.value !== oldPeriodicRefresh.value;
          const enabledIsChanged = periodicRefresh.enabled !== oldPeriodicRefresh.enabled;

          if (valueIsChanged || enabledIsChanged) {
            if (this.periodicRefreshInterval) {
              clearInterval(this.periodicRefreshInterval);
            }

            const periodicRefreshSeconds = toSeconds(periodicRefresh.value, periodicRefresh.unit);

            if (!periodicRefreshSeconds) {
              return;
            }

            this.periodicRefreshInterval = setInterval(() => {
              this.fetchList();
            }, periodicRefreshSeconds * 1000);
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
