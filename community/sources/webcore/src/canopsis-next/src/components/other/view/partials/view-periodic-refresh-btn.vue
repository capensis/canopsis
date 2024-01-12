<template>
  <v-tooltip top>
    <template #activator="{ on }">
      <v-btn
        v-on="on"
        :input-value="isPeriodicRefreshEnabled"
        :loading="pending"
        color="secondary"
        fab
        dark
        @click.stop="refreshHandler"
      >
        <v-icon v-if="!isPeriodicRefreshEnabled">
          refresh
        </v-icon>
        <v-progress-circular
          class="periodic-refresh-progress"
          v-else
          :rotate="270"
          :size="30"
          :width="2"
          :value="periodicRefreshProgressValue"
          color="white"
          button
        >
          <span class="refresh-btn">{{ periodicRefreshProgress | maxDurationByUnit }}</span>
        </v-progress-circular>
      </v-btn>
    </template>
    <span>{{ tooltipContent }}</span>
  </v-tooltip>
</template>

<script>
import { DATETIME_FORMATS } from '@/constants';

import { uid } from '@/helpers/uid';
import { convertDurationToString, durationToSeconds } from '@/helpers/date/duration';

import { activeViewMixin } from '@/mixins/active-view';
import { layoutNavigationEditingModeMixin } from '@/mixins/layout/navigation/editing-mode';

export default {
  inject: ['$periodicRefresh'],
  mixins: [
    activeViewMixin,
    layoutNavigationEditingModeMixin,
  ],
  data() {
    return {
      periodicRefreshInterval: null,
      periodicRefreshProgress: undefined,
    };
  },
  computed: {
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
      return this.periodicRefreshProgress / (this.periodicRefreshDelay / 100);
    },

    isPeriodicRefreshEnabled() {
      return this.view?.periodic_refresh?.enabled ?? false;
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
</script>

<style lang="scss">
.refresh-btn {
  text-decoration: none;
  text-transform: none;
}
</style>
