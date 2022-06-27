<template lang="pug">
  g(@dblclick="$emit('dblclick', $event)")
    rhombus-figure(
      :x="rhombus.x",
      :y="rhombus.y",
      :width="rhombus.width",
      :height="rhombus.height",
      :pointer-events="pointerEvents",
      fill="transparent",
      cursor="move",
      @mousedown.stop="$listeners.mousedown",
      @mouseup="$listeners.mouseup"
    )
    rect-selection(
      v-if="selected",
      :x="rhombus.x",
      :y="rhombus.y",
      :width="rhombus.width",
      :height="rhombus.height",
      :padding="padding",
      :color="color",
      :corner-radius="cornerRadius",
      resizable,
      @update="$listeners.update"
    )
    rect-connectors(
      v-if="connecting",
      :x="rhombus.x",
      :y="rhombus.y",
      :width="rhombus.width",
      :height="rhombus.height",
      :padding="padding",
      :color="color",
      @connected="$listeners.connected",
      @connecting="$listeners.connecting",
      @unconnect="$listeners.unconnect"
    )
</template>

<script>
import RectSelection from '../common/rect-selection.vue';
import RhombusFigure from '../common/rhombus-figure.vue';
import RectConnectors from '../common/rect-connectors.vue';

export default {
  inject: ['$mouseMove', '$mouseUp'],
  components: { RhombusFigure, RectSelection, RectConnectors },
  props: {
    rhombus: {
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
};
</script>
