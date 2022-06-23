<template lang="pug">
  g(@dblclick="$emit('dblclick', $event)")
    rect(
      :x="square.x",
      :y="square.y",
      :width="square.size",
      :height="square.size",
      :pointer-events="pointerEvents",
      fill="transparent",
      cursor="move",
      @mousedown.stop="$listeners.mousedown",
      @mouseup="$listeners.mouseup"
    )
    rect-selection(
      v-if="selected",
      :x="square.x",
      :y="square.y",
      :width="square.size",
      :height="square.size",
      :padding="padding",
      :color="color",
      :corner-radius="cornerRadius",
      resizable,
      @start:resize="startResize"
    )
    rect-connectors(
      v-if="connecting",
      :x="square.x",
      :y="square.y",
      :width="square.size",
      :height="square.size",
      :padding="padding",
      :color="color",
      @connected="$listeners.connected",
      @connecting="$listeners.connecting",
      @unconnect="$listeners.unconnect"
    )
</template>

<script>
import { resizeSquareShapeByDirection } from '../../utils/resize';

import RectSelection from '../common/rect-selection.vue';
import RectConnectors from '../common/rect-connectors.vue';

export default {
  inject: ['$mouseMove', '$mouseUp'],
  components: { RectSelection, RectConnectors },
  props: {
    square: {
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
      const { square, direction } = resizeSquareShapeByDirection(
        this.square,
        cursor,
        this.direction,
      );

      this.direction = direction;
      this.$emit('resize', square);
    },
  },
};
</script>
