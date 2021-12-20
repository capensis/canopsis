import { get } from 'lodash';

import { DATETIME_FORMATS } from '@/constants';

import Observer from '@/services/observer';

import uid from '@/helpers/uid';
import { convertDurationToString, durationToSeconds } from '@/helpers/date/duration';

import layoutNavigationEditingModeMixin from '../layout/navigation/editing-mode';

export default {
  mixins: [layoutNavigationEditingModeMixin],
  provide() {
    return {
      $periodicRefresh: this.$periodicRefresh,
    };
  },
  data() {
    return {
      periodicRefreshInterval: null,
      periodicRefreshProgress: undefined,
    };
  },

  beforeCreate() {
    this.$periodicRefresh = new Observer();
  },

  watch: {
    isPeriodicRefreshEnabled(value, oldValue) {
      if (value && (!oldValue || !this.periodicRefreshInterval) && !this.isNavigationEditingMode) {
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

    isNavigationEditingMode(value, oldValue) {
      if (value !== oldValue && this.isPeriodicRefreshEnabled) {
        if (this.popupId) {
          this.$popups.remove({ id: this.popupId });
        }

        if (value) {
          this.popupId = uid('popup');

          this.$popups.info({
            id: this.popupId,
            text: this.$t('layout.sideBar.ordering.popups.periodicRefreshWasPaused'),
            autoClose: 7000,
          });

          this.stopPeriodicRefreshInterval();
        } else {
          this.popupId = uid('popup');

          this.$popups.info({
            id: this.popupId,
            text: this.$t('layout.sideBar.ordering.popups.periodicRefreshWasResumed'),
          });

          this.resumePeriodicRefreshInterval();
        }
      }
    },
  },

  mounted() {
    if (this.isPeriodicRefreshEnabled && !this.isNavigationEditingMode) {
      this.startPeriodicRefreshInterval();
    }
  },

  beforeDestroy() {
    this.stopPeriodicRefreshInterval();
  },

  computed: {
    periodicRefreshProgressFormatted() {
      return convertDurationToString(
        this.periodicRefreshProgress,
        DATETIME_FORMATS.refreshFieldFormat,
      );
    },

    periodicRefreshProgressValue() {
      return this.periodicRefreshProgress / (this.periodicRefreshDelay / 100);
    },

    isPeriodicRefreshEnabled() {
      return get(this.view, 'periodic_refresh.enabled', false);
    },

    periodicRefreshDelay() {
      return this.view?.periodic_refresh
        ? durationToSeconds(this.view.periodic_refresh)
        : 0;
    },

    refreshHandler() {
      return this.isPeriodicRefreshEnabled && !this.isNavigationEditingMode
        ? this.callSubscribers
        : this.refresh;
    },
  },
  methods: {
    refresh() {
      return this.$periodicRefresh.notify();
    },

    async callSubscribers() {
      this.stopPeriodicRefreshInterval();

      await this.refresh();

      this.startPeriodicRefreshInterval();
    },

    resetRefreshInterval() {
      this.periodicRefreshProgress = this.periodicRefreshDelay;
    },

    refreshTick() {
      if (this.periodicRefreshProgress <= 0) {
        this.callSubscribers();
      } else {
        this.periodicRefreshProgress -= 1;
      }
    },

    resumePeriodicRefreshInterval() {
      if (this.periodicRefreshInterval) {
        this.stopPeriodicRefreshInterval();
      }

      this.periodicRefreshInterval = setInterval(this.refreshTick, 1000);
    },

    startPeriodicRefreshInterval() {
      this.resetRefreshInterval();
      this.resumePeriodicRefreshInterval();
    },

    stopPeriodicRefreshInterval() {
      clearInterval(this.periodicRefreshInterval);

      this.periodicRefreshInterval = undefined;
    },
  },
};
