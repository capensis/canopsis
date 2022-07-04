<template lang="pug">
  g(@dblclick="$emit('dblclick', $event)")
    document-figure(
      :x="document.x",
      :y="document.y",
      :width="document.width",
      :height="document.height",
      :offset="document.offset",
      :pointer-events="pointerEvents",
      fill="transparent",
      cursor="move",
      @mousedown.stop="$listeners.mousedown",
      @mouseup="$listeners.mouseup"
    )
    rect-selection(
      v-if="selected",
      :x="document.x",
      :y="document.y",
      :width="document.width",
      :height="document.height",
      :padding="padding",
      :color="color",
      :corner-radius="cornerRadius",
      resizable,
      @update="$listeners.update"
    )
    rect-connectors(
      v-if="connecting",
      :x="document.x",
      :y="document.y",
      :width="document.width",
      :height="document.height",
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
import DocumentFigure from '../common/document-figure.vue';

export default {
  components: { DocumentFigure, RectSelection, RectConnectors },
  props: {
    document: {
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
