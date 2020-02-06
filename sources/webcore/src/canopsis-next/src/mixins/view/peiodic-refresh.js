import { get } from 'lodash';
import { DATETIME_FORMATS } from '@/constants';
import { getSecondsByUnit } from '@/helpers/time';

export default {
  data() {
    return {
      periodicRefreshInterval: null,
      periodicRefreshProgress: undefined,
    };
  },

  watch: {
    isPeriodicRefreshEnabled(value, oldValue) {
      if (value && (!oldValue || !this.periodicRefreshInterval)) {
        this.startPeriodicRefreshInterval();
      } else if (oldValue && !value) {
        this.stopPeriodicRefreshInterval();
      }
    },
    periodicRefreshDelay(value, oldValue) {
      if (value !== oldValue) {
        this.resetRefreshInterval();
      }
    },
  },


  mounted() {
    if (this.isPeriodicRefreshEnabled) {
      this.startPeriodicRefreshInterval();
    }
  },


  beforeDestroy() {
    this.stopPeriodicRefreshInterval();
  },

  computed: {
    periodicRefreshProgressFormated() {
      return this.$options.filters.duration(
        this.periodicRefreshProgress,
        undefined,
        DATETIME_FORMATS.refreshFieldFormat,
      );
    },

    periodicRefreshProgressValue() {
      return this.periodicRefreshProgress / (this.periodicRefreshDelay / 100);
    },

    isPeriodicRefreshEnabled() {
      return get(this.view, 'periodicRefresh.enabled', false);
    },

    periodicRefreshUnit() {
      return get(this.view, 'periodicRefresh.unit');
    },

    periodicRefreshValue() {
      return get(this.view, 'periodicRefresh.interval') || get(this.view, 'periodicRefresh.value', 0);
    },

    periodicRefreshDelay() {
      return getSecondsByUnit(this.periodicRefreshValue, this.periodicRefreshUnit);
    },

    refreshHandler() {
      return this.isPeriodicRefreshEnabled ? this.refreshViewWithProgress : this.refreshView;
    },
  },
  methods: {
    async refreshView() {
      await this.fetchView({ id: this.id });

      if (this.activeTab) {
        this.forceUpdateQuery({ id: this.activeTab._id });
      }
    },

    async refreshViewWithProgress() {
      this.stopPeriodicRefreshInterval();

      await this.refreshView();

      this.startPeriodicRefreshInterval();
    },

    resetRefreshInterval() {
      this.periodicRefreshProgress = this.periodicRefreshDelay;
    },

    refreshTick() {
      if (this.periodicRefreshProgress <= 0) {
        this.refreshViewWithProgress();
      } else {
        this.periodicRefreshProgress -= 1;
      }
    },

    startPeriodicRefreshInterval() {
      this.resetRefreshInterval();

      if (this.periodicRefreshInterval) {
        this.stopPeriodicRefreshInterval();
      }

      this.periodicRefreshInterval = setInterval(this.refreshTick, 1000);
    },

    stopPeriodicRefreshInterval() {
      clearInterval(this.periodicRefreshInterval);

      this.periodicRefreshInterval = undefined;
    },
  },
};
