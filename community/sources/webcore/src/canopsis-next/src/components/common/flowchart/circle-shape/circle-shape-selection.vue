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
      :x="circle.x",
      :y="circle.y",
      :width="circle.diameter",
      :height="circle.diameter",
      :padding="padding",
      :color="color",
      :corner-radius="cornerRadius",
      resizable,
      aspect-ratio,
      @update="onRectUpdate"
    )
    rect-connectors(
      v-if="connecting",
      :x="circle.x",
      :y="circle.y",
      :width="circle.diameter",
      :height="circle.diameter",
      :padding="padding",
      :color="color",
      @connected="$listeners.connected",
      @connecting="$listeners.connecting",
      @unconnect="$listeners.unconnect"
    )
</template>

<script>
import RectSelection from '../common/rect-selection.vue';
import CircleFigure from '../common/circle-figure.vue';
import RectConnectors from '../common/rect-connectors.vue';

export default {
  components: { RectConnectors, CircleFigure, RectSelection },
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
    connecting: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    onRectUpdate(rect) {
      this.$emit('update', {
        diameter: rect.width,
        x: rect.x,
        y: rect.y,
      });
    },
  },
};
</script>
