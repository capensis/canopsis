<template lang="pug">
  v-navigation-drawer(
    v-model="isOpen",
    :ignore-click-outside="hasMaximizedModal",
    :custom-close-conditional="closeCondition",
    :width="400",
    right,
    fixed,
    temporary
  )
    div(v-if="title")
      v-toolbar(color="secondary")
        v-list
          v-list-tile
            v-list-tile-title.white--text {{ title }}
        v-btn(icon, @click.stop="closeHandler")
          v-icon(color="white") close
      v-divider
      // @slot use this slot default
    slot
</template>

<script>
import ClickOutside from '@/services/click-outside';

import { sideBarInnerMixin } from '@/mixins/side-bar/side-bar-inner';

/**
 * Wrapper for each modal window
 *
 * @prop {Object} [navigationDrawerProps={}] - Properties for vuetify v-navigation-drawer
 */
export default {
  provide() {
    return {
      $clickOutside: this.$clickOutside,
    };
  },
  mixins: [sideBarInnerMixin],
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
  beforeCreate() {
    this.$clickOutside = new ClickOutside();
  },
  methods: {
    closeHandler() {
      if (this.closeCondition()) {
        this.hideSideBar();
      }
    },

    closeCondition(...args) {
      return this.$clickOutside.call(...args);
    },
  },
};
</script>
