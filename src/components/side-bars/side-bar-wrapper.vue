<template lang="pug">
  v-navigation-drawer(
  v-model="isOpen",
  :stateless="hasModals",
  v-bind="navigationDrawerProps"
  )
    div(v-if="title")
      v-toolbar(color="blue darken-4")
        v-list
          v-list-tile
            v-list-tile-title.white--text.text-xs-center {{ title }}
        v-btn(@click.stop="hideSideBar", icon)
          v-icon(color="white") close
      v-divider
      // @slot use this slot default
    slot
</template>

<script>
import settingsInnerMixin from '@/mixins/side-bar/side-bar-inner';
import { SIDE_BARS } from '@/constants';

/**
 * Wrapper for each modal window
 *
 * @prop {Object} modal - The current modal object
 * @prop {number} index - The current modal index in the store
 * @prop {Object} [dialogProps={}] - Properties for vuetify v-dialog
 */
export default {
  mixins: [settingsInnerMixin],
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
