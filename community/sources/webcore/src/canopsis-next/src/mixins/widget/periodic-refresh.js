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
    periodicRefreshSettings() {
      return { ...this.widget.parameters.periodic_refresh };
    },

    periodicRefreshEnabled() {
      return this.widget.parameters.periodic_refresh?.enabled;
    },

    periodicRefreshSeconds() {
      const { value, unit } = this.widget.parameters.periodic_refresh ?? {};

      return toSeconds(value, unit);
    },

    isPeriodicRefreshVisible() {
      return this.visible !== false;
    },
  },
  watch: {
    periodicRefreshSettings: {
      immediate: true,
      handler(value, oldValue) {
        const periodicRefresh = value;
        const oldPeriodicRefresh = oldValue ?? {};
        const valueIsChanged = periodicRefresh?.value !== oldPeriodicRefresh.value;
        const unitIsChanged = periodicRefresh?.unit !== oldPeriodicRefresh.unit;

        if (valueIsChanged || unitIsChanged) {
          this.resetPeriodicRefresh();
        }

        if (this.activeViewPeriodicRefreshPaused || !this.isPeriodicRefreshVisible) {
          return;
        }

        if (periodicRefresh?.enabled && periodicRefresh?.value) {
          const enabledIsChanged = periodicRefresh.enabled !== oldPeriodicRefresh.enabled;

          if (enabledIsChanged) {
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
      if (!this.isPeriodicRefreshVisible) {
        return;
      }

      if (value) {
        this.pausePriodicRefresh();

        return;
      }

      this.resumePriodicRefresh();
    },

    visible(value) {
      if (value === false) {
        this.stopPeriodicRefresh();

        return;
      }

      if (
        this.activeViewPeriodicRefreshPaused ||
        !this.periodicRefreshEnabled ||
        !this.periodicRefreshSeconds
      ) {
        return;
      }

      this.startPeriodicRefresh();
    },
  },
  beforeUnmount() {
    this.stopPeriodicRefresh();
  },
  methods: {
    resetPeriodicRefresh() {
      this.intervalDelay = this.periodicRefreshSeconds * 1000;

      if (
        this.periodicRefreshInterval &&
        !this.activeViewPeriodicRefreshPaused &&
        this.isPeriodicRefreshVisible
      ) {
        this.startPeriodicRefresh();
      }
    },

    startPeriodicRefresh() {
      if (
        !this.isPeriodicRefreshVisible ||
        !this.periodicRefreshEnabled ||
        !this.periodicRefreshSeconds
      ) {
        return;
      }

      if (this.periodicRefreshInterval) {
        this.stopPeriodicRefresh();
      }

      this.periodicRefreshInterval = setInterval(
        this.fetchList,
        this.periodicRefreshSeconds * 1000,
      );

      this.intervalStartedAt = Date.now();
    },

    resumePriodicRefresh() {
      if (
        !this.isPeriodicRefreshVisible ||
        !this.periodicRefreshEnabled ||
        !this.periodicRefreshSeconds
      ) {
        return;
      }

      setTimeout(() => {
        this.fetchList();
        this.startPeriodicRefresh();
      }, this.intervalDelay || 0);
    },

    pausePriodicRefresh() {
      if (!this.periodicRefreshInterval) {
        return;
      }

      clearInterval(this.periodicRefreshInterval);

      this.intervalDelay =
        this.periodicRefreshSeconds * 1000 - (Date.now() - this.intervalStartedAt);
    },

    stopPeriodicRefresh() {
      clearInterval(this.periodicRefreshInterval);

      this.periodicRefreshInterval = null;
      this.intervalStartedAt = null;
      this.intervalDelay = null;
    },
  },
};
