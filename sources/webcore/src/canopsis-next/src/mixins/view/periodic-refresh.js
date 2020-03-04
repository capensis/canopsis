import { get } from 'lodash';

import { DATETIME_FORMATS } from '@/constants';

import uid from '@/helpers/uid';
import { getSecondsByUnit } from '@/helpers/time';
import subscribersMixin from '@/mixins/subscribers';

import layoutNavigationEditingModeMixin from '../layout/navigation/editing-mode';

export default {
  mixins: [subscribersMixin, layoutNavigationEditingModeMixin],
  provide() {
    return {
      $periodicRefresh: {
        subscribe: this.subscribe,
        unsubscribe: this.unsubscribe,
        subscribers: this.subscribers,
      },
    };
  },
  data() {
    return {
      periodicRefreshInterval: null,
      periodicRefreshProgress: undefined,
    };
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
      return this.isPeriodicRefreshEnabled && !this.isNavigationEditingMode ?
        this.callSubscribers : this.refreshView;
    },
  },
  methods: {
    async callSubscribers() {
      this.stopPeriodicRefreshInterval();

      await this.notify();

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
