<template lang="pug">
  g(@dblclick="$emit('dblclick', $event)")
    circle-figure(
      :x="circle.x",
      :y="circle.y",
      :diameter="circle.diameter",
      :pointer-events="pointerEvents",
      fill="transparent",
      cursor="move",
      pointer-events="all",
      @mousedown.stop="$listeners.mousedown",
      @mouseup="$listeners.mouseup"
    )
    rect-selection(
      v-if="selected",
      :rect="rect",
      :padding="padding",
      :color="color",
      :corner-radius="cornerRadius",
      resizable,
      @start:resize="startResize"
    )
</template>

<script>
import { resizeSquareShapeByDirection } from '../../utils/resize';

import RectSelection from '../common/rect-selection.vue';
import CircleFigure from '../common/circle-figure.vue';

export default {
  inject: ['$mouseMove', '$mouseUp'],
  components: { CircleFigure, RectSelection },
  props: {
    circle: {
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
  },
  data() {
    return {
      direction: undefined,
    };
  },
  computed: {
    rect() {
      return {
        x: this.circle.x,
        y: this.circle.y,
        width: this.circle.diameter,
        height: this.circle.diameter,
      };
    },
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
        {
          x: this.circle.x,
          y: this.circle.y,
          size: this.circle.diameter,
        },
        cursor,
        this.direction,
      );

      this.direction = direction;
      this.$emit('resize', {
        diameter: square.size,
        x: square.x,
        y: square.y,
      });
    },
  },
};
</script>
