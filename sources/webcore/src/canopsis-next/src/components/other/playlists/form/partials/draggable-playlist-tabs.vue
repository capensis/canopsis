<template lang="pug">
  draggable.tabs-draggable-panel.secondary.lighten-1(
    :value="tabs",
    :class="{ empty: isTabsEmpty, disabled: disabled }",
    :options="draggableOptions",
    @change="changeTabsOrdering"
  )
    tab-panel-content(v-for="tab in tabs", :tab="tab", hideActions, :key="tab._id")
      playlist-tab-item(slot="title", :tab="tab")
</template>

<script>
import Draggable from 'vuedraggable';

import { VUETIFY_ANIMATION_DELAY } from '@/config';

import TabPanelContent from '@/components/other/playlists/form/partials/tab-panel-content.vue';
import PlaylistTabItem from '@/components/other/playlists/form/partials/playlist-tab-item.vue';
import { dragDropChangePositionHandler } from '@/helpers/dragdrop';

export default {
  components: { PlaylistTabItem, TabPanelContent, Draggable },
  model: {
    prop: 'tabs',
    event: 'change',
  },
  props: {
    tabs: {
      type: Array,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    isTabsEmpty() {
      return this.tabs.length === 0;
    },

    draggableOptions() {
      return {
        animation: VUETIFY_ANIMATION_DELAY,
        disabled: this.disabled,
      };
    },
  },
  methods: {
    changeTabsOrdering(event) {
      this.$emit('change', dragDropChangePositionHandler(this.tabs, event));
    },
  },
};
</script>

<style lang="scss" scoped>
  .tabs-draggable-panel {
    &:not(.disabled) /deep/ .tab-panel-item {
      cursor: move;
    }

    &.empty {
      &:after {
        content: '';
        display: block;
        height: 48px;
        border: 4px dashed #4f6479;
        border-radius: 5px;
        position: relative;
      }
    }
  }
</style>
