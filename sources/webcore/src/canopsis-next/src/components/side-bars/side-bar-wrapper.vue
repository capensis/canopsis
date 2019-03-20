<template lang="pug">
  v-navigation-drawer(
  v-model="isOpen",
  :stateless="hasModals",
  v-bind="navigationDrawerProps",
  )
    div(v-if="title")
      v-toolbar(color="secondary")
        v-list
          v-list-tile
            v-list-tile-title.white--text {{ title }}
        v-btn(@click.stop="hideSideBar", icon)
          v-icon(color="white") close
      v-divider
      // @slot use this slot default
    slot
</template>

<script>
import { SIDE_BARS } from '@/constants';

import sideBarInnerMixin from '@/mixins/side-bar/side-bar-inner';

/**
 * Wrapper for each modal window
 *
 * @prop {Object} [navigationDrawerProps={}] - Properties for vuetify v-navigation-drawer
 */
export default {
  mixins: [sideBarInnerMixin],
  props: {
    navigationDrawerProps: {
      type: Object,
      default: () => ({}),
    },
  },
  data() {
    return {
      ready: false,
    };
  },
  computed: {
    title() {
      const TITLES_MAP = {
        [SIDE_BARS.alarmSettings]: this.$t('settings.titles.alarmListSettings'),
        [SIDE_BARS.contextSettings]: this.$t('settings.titles.contextTableSettings'),
        [SIDE_BARS.weatherSettings]: this.$t('settings.titles.weatherSettings'),
        [SIDE_BARS.statsHistogramSettings]: this.$t('settings.titles.statsHistogramSettings'),
        [SIDE_BARS.statsCurvesSettings]: this.$t('settings.titles.statsCurvesSettings'),
        [SIDE_BARS.statsTableSettings]: this.$t('settings.titles.statsTableSettings'),
        [SIDE_BARS.statsCalendarSettings]: this.$t('settings.titles.statsCalendarSettings'),
        [SIDE_BARS.statsNumberSettings]: this.$t('settings.titles.statsNumberSettings'),
        [SIDE_BARS.statsNumberSettings]: this.$t('settings.titles.statsNumberSettings'),
        [SIDE_BARS.statsTextSettings]: this.$t('settings.titles.statsTextSettings'),
      };

      return this.sideBarConfig.sideBarTitle || TITLES_MAP[this.sideBarName];
    },
    isOpen: {
      get() {
        return !this.isSideBarHidden && this.ready && this.sideBarName;
      },
      set(value) {
        if (!value) {
          this.hideSideBar();
        }
      },
    },
  },
  mounted() {
    this.ready = true;
  },
};
</script>
