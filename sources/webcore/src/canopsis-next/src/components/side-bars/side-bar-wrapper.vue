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
      const sideBars = { ...this.$constants.SIDE_BARS };
      const TITLES_MAP = {
        [sideBars.alarmSettings]: this.$t('settings.titles.alarmListSettings'),
        [sideBars.contextSettings]: this.$t('settings.titles.contextTableSettings'),
        [sideBars.weatherSettings]: this.$t('settings.titles.weatherSettings'),
        [sideBars.statsHistogramSettings]: this.$t('settings.titles.statsHistogramSettings'),
        [sideBars.statsCurvesSettings]: this.$t('settings.titles.statsCurvesSettings'),
        [sideBars.statsTableSettings]: this.$t('settings.titles.statsTableSettings'),
        [sideBars.statsCalendarSettings]: this.$t('settings.titles.statsCalendarSettings'),
        [sideBars.statsNumberSettings]: this.$t('settings.titles.statsNumberSettings'),
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
