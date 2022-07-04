<template lang="pug">
  g(@dblclick="$emit('dblclick', $event)")
    process-figure(
      :x="process.x",
      :y="process.y",
      :width="process.width",
      :height="process.height",
      :offset="process.offset",
      :pointer-events="pointerEvents",
      fill="transparent",
      cursor="move",
      @mousedown.stop="$listeners.mousedown",
      @mouseup="$listeners.mouseup"
    )
    rect-selection(
      v-if="selected",
      :x="process.x",
      :y="process.y",
      :width="process.width",
      :height="process.height",
      :padding="padding",
      :color="color",
      :corner-radius="cornerRadius",
      resizable,
      @update="$listeners.update"
    )
    rect-connectors(
      v-if="connecting",
      :x="process.x",
      :y="process.y",
      :width="process.width",
      :height="process.height",
      :padding="padding",
      :color="color",
      @connected="$listeners.connected",
      @connecting="$listeners.connecting",
      @unconnect="$listeners.unconnect"
    )
</template>

<script>
import RectSelection from '../common/rect-selection.vue';
import RectConnectors from '../common/rect-connectors.vue';
import ProcessFigure from '../common/process-figure.vue';

export default {
  components: { ProcessFigure, RectSelection, RectConnectors },
  props: {
    process: {
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
