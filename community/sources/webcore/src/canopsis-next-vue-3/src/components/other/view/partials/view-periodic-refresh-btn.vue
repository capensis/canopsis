<template>
  <v-tooltip location="top">
    <template #activator="{ props }">
      <v-btn
        v-bind="props"
        :active="isPeriodicRefreshEnabled"
        :loading="activeViewPending"
        color="secondary"
        size="x-large"
        icon
        @click.stop="refreshHandler"
      >
        <v-icon v-if="!isPeriodicRefreshEnabled" icon="refresh" />
        <v-progress-circular
          v-else
          :model-value="periodicRefreshProgressValue"
          :rotate="270"
          :size="40"
          :width="2"
          color="white"
        >
          <span class="refresh-btn">
            {{ periodicRefreshProgress }}
          </span>
        </v-progress-circular>
      </v-btn>
    </template>
    <span>{{ tooltipContent }}</span>
  </v-tooltip>
</template>

<script>
import { DATETIME_FORMATS } from '@/constants';

import { convertDurationToString, durationToSeconds } from '@/helpers/date/duration';
import { uid } from '@/helpers/uid';

import { activeViewMixin } from '@/mixins/active-view';
import { layoutNavigationEditingModeMixin } from '@/mixins/layout/navigation/editing-mode';

export default {
  inject: ['$periodicRefresh'],
  mixins: [activeViewMixin, layoutNavigationEditingModeMixin],
  data() {
    return {
      popupId: null,
      periodicRefreshInterval: null,
      periodicRefreshProgress: undefined,
    };
  },
  computed: {
    periodicRefreshFullPaused() {
      return this.activeViewPeriodicRefreshPaused || this.isNavigationEditingMode;
    },

    tooltipContent() {
      return this.isPeriodicRefreshEnabled
        ? this.periodicRefreshProgressFormatted
        : this.$t('common.refresh');
    },

    periodicRefreshProgressFormatted() {
      return convertDurationToString(
        this.periodicRefreshProgress,
        DATETIME_FORMATS.refreshFieldFormat,
      );
    },

    periodicRefreshProgressValue() {
      return this.periodicRefreshDelay
        ? this.periodicRefreshProgress / (this.periodicRefreshDelay / 100)
        : 0;
    },

    isPeriodicRefreshEnabled() {
      return this.view?.periodic_refresh?.enabled ?? false;
    },

    periodicRefreshDelay() {
      return this.view?.periodic_refresh ? durationToSeconds(this.view.periodic_refresh) : 0;
    },

    refreshHandler() {
      return this.isPeriodicRefreshEnabled && !this.periodicRefreshFullPaused
        ? this.callSubscribers
        : this.refresh;
    },
  },
  watch: {
    isPeriodicRefreshEnabled(value, oldValue) {
      if (
        value &&
        (!oldValue || !this.periodicRefreshInterval) &&
        !this.periodicRefreshFullPaused
      ) {
        this.startPeriodicRefreshInterval();
      } else if (oldValue && !value) {
        this.stopPeriodicRefreshInterval();
      } else if (value && !oldValue) {
        this.resetRefreshInterval();
      }
    },

    periodicRefreshDelay(value, oldValue) {
      if (value !== oldValue) {
        this.resetRefreshInterval();
      }
    },

    isNavigationEditingMode(value) {
      this.periodicRefreshPausedWatcher(
        value,
        this.$t('layout.sideBar.ordering.popups.periodicRefreshWasPausedWhileEditingGroups'),
      );
    },

    activeViewPeriodicRefreshPaused(value) {
      this.periodicRefreshPausedWatcher(
        value,
        this.$t('layout.sideBar.ordering.popups.periodicRefreshWasPaused'),
      );
    },
  },

  mounted() {
    if (this.isPeriodicRefreshEnabled && !this.periodicRefreshFullPaused) {
      this.startPeriodicRefreshInterval();
    }
  },

  beforeUnmount() {
    this.stopPeriodicRefreshInterval();
  },

  methods: {
    periodicRefreshPausedWatcher(value, pausedPopup) {
      if (
        value !== Boolean(this.periodicRefreshInterval) ||
        value !== this.periodicRefreshFullPaused
      ) {
        return;
      }

      if (this.popupId) {
        this.$popups.remove({ id: this.popupId });
        this.popupId = null;
      }

      if (value && this.periodicRefreshFullPaused) {
        this.stopPeriodicRefreshInterval();

        this.popupId = uid('popup');
        this.$popups.info({
          id: this.popupId,
          text: pausedPopup,
          autoClose: 7000,
        });

        return;
      }

      this.resumePeriodicRefreshInterval();

      this.popupId = uid('popup');
      this.$popups.info({
        id: this.popupId,
        text: this.$t('layout.sideBar.ordering.popups.periodicRefreshWasResumed'),
      });
    },

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
</script>

<style lang="scss" scoped>
.refresh-btn {
  text-decoration: none;
  text-transform: none;
}
</style>
