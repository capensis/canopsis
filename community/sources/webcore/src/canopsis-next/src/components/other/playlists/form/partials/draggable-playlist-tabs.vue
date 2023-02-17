<template lang="pug">
  v-layout(column)
    v-layout
      v-flex.text-xs-center.font-weight-bold(xs4) {{ $tc('common.group') }}
      v-flex.text-xs-center.font-weight-bold(xs4) {{ $tc('common.view') }}
      v-flex.text-xs-center.font-weight-bold(xs4) {{ $tc('common.tab') }}
    v-layout(column)
      draggable.tabs-draggable-panel.secondary.lighten-1(
        :value="tabs",
        :class="{ empty: isTabsEmpty, disabled: disabled }",
        :options="draggableOptions",
        @change="changeTabsOrdering"
      )
        tab-panel-content(v-for="tab in tabs", :tab="tab", hideActions, :key="tab._id")
          template(#title="")
            playlist-tab-item(:tab="tab")
</template>

<script>
import Draggable from 'vuedraggable';

import { VUETIFY_ANIMATION_DELAY } from '@/config';

import { dragDropChangePositionHandler } from '@/helpers/dragdrop';

import TabPanelContent from './tab-panel-content.vue';
import PlaylistTabItem from './playlist-tab-item.vue';

export default {
  components: { Draggable, TabPanelContent, PlaylistTabItem },
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
    &:not(.disabled) ::v-deep .tab-panel-item {
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
