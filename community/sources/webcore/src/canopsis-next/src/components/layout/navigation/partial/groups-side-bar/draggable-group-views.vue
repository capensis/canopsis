<template lang="pug">
  draggable.views-panel.secondary.lighten-1(
    :value="views",
    :class="{ empty: isViewsEmpty }",
    :options="draggableOptions",
    @change="changeViewsOrdering"
  )
    group-view-panel(
      v-for="view in views",
      :key="view._id",
      :view="view"
    )
</template>

<script>
import Draggable from 'vuedraggable';

import { VUETIFY_ANIMATION_DELAY } from '@/config';

import { dragDropChangePositionHandler } from '@/helpers/dragdrop';

import GroupViewPanel from './group-view-panel.vue';

export default {
  components: { GroupViewPanel, Draggable },
  model: {
    prop: 'views',
    event: 'change',
  },
  props: {
    views: {
      type: Array,
      required: true,
    },
    put: {
      type: Boolean,
      default: false,
    },
    pull: {
      type: [Boolean, String],
      default: false,
    },
  },
  computed: {
    isViewsEmpty() {
      return this.views.length === 0;
    },

    draggableOptions() {
      return {
        animation: VUETIFY_ANIMATION_DELAY,
        group: { name: 'views', put: this.put, pull: this.pull },
      };
    },
  },
  methods: {
    changeViewsOrdering(event) {
      this.$emit('change', dragDropChangePositionHandler(this.views, event));
    },
  },
};
</script>

<style lang="scss" scoped>
  .views-panel {
    & ::v-deep .panel-item-content {
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
