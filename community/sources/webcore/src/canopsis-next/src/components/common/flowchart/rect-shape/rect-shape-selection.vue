<template lang="pug">
  g(@dblclick="$emit('dblclick', $event)")
    rect(
      :x="rect.x",
      :y="rect.y",
      :width="rect.width",
      :height="rect.height",
      :pointer-events="pointerEvents",
      fill="transparent",
      cursor="move",
      @mousedown.stop="$listeners.mousedown",
      @mouseup="$listeners.mouseup"
    )
    rect-selection(
      v-if="selected",
      :x="rect.x",
      :y="rect.y",
      :width="rect.width",
      :height="rect.height",
      :aspect-ratio="rect.aspectRatio",
      :padding="padding",
      :color="color",
      :corner-radius="cornerRadius",
      cursor="move",
      @update="$listeners.update"
    )
    rect-connectors(
      v-if="connecting",
      :x="rect.x",
      :y="rect.y",
      :width="rect.width",
      :height="rect.height",
      :color="color",
      @connected="$listeners.connected",
      @connecting="$listeners.connecting",
      @unconnect="$listeners.unconnect"
    )
</template>

<script>
import RectSelection from '../common/rect-selection.vue';
import RectConnectors from '../common/rect-connectors.vue';

export default {
  components: { RectConnectors, RectSelection },
  props: {
    rect: {
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
