<template lang="pug">
  v-navigation-drawer(
    data-test="sideBarWrapper",
    v-model="isOpen",
    :ignore-click-outside="hasMaximizedModal",
    v-bind="navigationDrawerProps"
  )
    div(v-if="title")
      v-toolbar(color="secondary")
        v-list
          v-list-tile
            v-list-tile-title.white--text {{ title }}
        v-btn(data-test="closeWidget", @click.stop="hideSideBar", icon)
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
      if (this.sideBarConfig.sideBarTitle) {
        return this.sideBarConfig.sideBarTitle;
      }

      return this.sideBarName ? this.$t(`settings.titles.${this.sideBarName}`) : '';
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
