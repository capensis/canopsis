import { toSeconds } from '@/helpers/date/duration';

import { activeViewMixin } from '@/mixins/active-view';

export const widgetPeriodicRefreshMixin = {
  mixins: [activeViewMixin],

  data() {
    return {
      periodicRefreshInterval: null,
    };
  },
  computed: {
    periodicRefreshEnabled() {
      return this.widget.parameters.periodic_refresh?.enabled;
    },

    periodicRefreshSeconds() {
      const { value, unit } = this.widget.parameters.periodic_refresh ?? {};

      return toSeconds(value, unit);
    },
  },
  watch: {
    'widget.parameters.periodic_refresh': {
      immediate: true,
      handler(value, oldValue) {
        const periodicRefresh = value;
        const oldPeriodicRefresh = oldValue ?? {};

        if (this.activeViewPeriodicRefreshPaused) {
          return;
        }

        if (periodicRefresh?.enabled && periodicRefresh?.value) {
          const valueIsChanged = periodicRefresh.value !== oldPeriodicRefresh.value;
          const enabledIsChanged = periodicRefresh.enabled !== oldPeriodicRefresh.enabled;

          if (valueIsChanged || enabledIsChanged) {
            if (this.periodicRefreshInterval) {
              this.stopPeriodicRefresh();
            }

            if (!this.periodicRefreshSeconds) {
              return;
            }

            this.startPeriodicRefresh();
          }
        } else {
          this.stopPeriodicRefresh();
        }
      },
    },

    activeViewPeriodicRefreshPaused(value) {
      if (value) {
        this.pausePriodicRefresh();

        return;
      }

      this.resumePriodicRefresh();
    },
  },
  beforeDestroy() {
    this.stopPeriodicRefresh();
  },
  methods: {
    startPeriodicRefresh() {
      if (this.periodicRefreshInterval) {
        this.stopPeriodicRefresh();
      }

      this.periodicRefreshInterval = setInterval(this.fetchList, this.periodicRefreshSeconds * 1000);

      this.intervalStartedAt = Date.now();
    },

    resumePriodicRefresh() {
      setTimeout(() => {
        this.fetchList();
        this.startPeriodicRefresh();
      }, this.intervalDelay || 0);
    },

    pausePriodicRefresh() {
      clearInterval(this.periodicRefreshInterval);

      this.intervalDelay = (this.periodicRefreshSeconds * 1000) - (Date.now() - this.intervalStartedAt);
    },

    stopPeriodicRefresh() {
      clearInterval(this.periodicRefreshInterval);

      this.periodicRefreshInterval = null;
      this.intervalStartedAt = null;
      this.intervalDelay = null;
    },
  },
};
