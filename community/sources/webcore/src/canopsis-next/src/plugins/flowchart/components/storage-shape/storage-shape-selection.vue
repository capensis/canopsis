<template lang="pug">
  g(@dblclick="$emit('dblclick', $event)")
    storage-figure(
      :x="storage.x",
      :y="storage.y",
      :width="storage.width",
      :height="storage.height",
      :radius="storage.radius",
      :pointer-events="pointerEvents",
      fill="transparent",
      cursor="move",
      @mousedown.stop="$listeners.mousedown",
      @mouseup="$listeners.mouseup"
    )
    rect-selection(
      v-if="selected",
      :x="storage.x",
      :y="storage.y",
      :width="storage.width",
      :height="storage.height",
      :padding="padding",
      :color="color",
      :corner-radius="cornerRadius",
      resizable,
      @start:resize="startResize"
    )
    rect-connectors(
      v-if="connecting",
      :x="storage.x",
      :y="storage.y",
      :width="storage.width",
      :height="storage.height",
      :padding="padding",
      :color="color",
      @connected="$listeners.connected",
      @connecting="$listeners.connecting",
      @unconnect="$listeners.unconnect"
    )
</template>

<script>
import { resizeRectangleShapeByDirection } from '../../utils/resize';

import RectSelection from '../common/rect-selection.vue';
import StorageFigure from '../common/storage-figure.vue';
import RectConnectors from '../common/rect-connectors.vue';

export default {
  inject: ['$mouseMove', '$mouseUp'],
  components: { StorageFigure, RectSelection, RectConnectors },
  props: {
    storage: {
      type: Object,
      required: true,
    },
    selected: {
      type: Boolean,
      default: false,
    },
    padding: {
      type: Number,
      default: 0,
    },
    color: {
      type: String,
      default: 'blue',
    },
    pointerEvents: {
      type: String,
      default: 'all',
    },
    cornerRadius: {
      type: Number,
      default: 4,
    },
    connecting: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      direction: undefined,
    };
  },
  methods: {
    startResize(direction) {
      this.$mouseMove.register(this.onResize);
      this.$mouseUp.register(this.finishResize);
      this.direction = direction;
    },

    finishResize() {
      this.$mouseMove.unregister(this.onResize);
      this.$mouseUp.unregister(this.finishResize);
    },

    onResize(cursor) {
      const { rect, direction } = resizeRectangleShapeByDirection(
        this.storage,
        cursor,
        this.direction,
      );

      this.direction = direction;
      this.$emit('resize', rect);
    },
  },
};
</script>
