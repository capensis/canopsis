<template>
  <g @click.stop="">
    <rect-selection
      v-if="selected"
      :x="shape.x"
      :y="shape.y"
      :width="shape.width"
      :height="shape.height"
      :aspect-ratio="shape.aspectRatio"
      :padding="padding"
      :color="color"
      :corner-radius="cornerRadius"
      cursor="move"
      @update="$listeners.update"
    />
    <rect-connectors
      v-if="connecting"
      :x="shape.x"
      :y="shape.y"
      :width="shape.width"
      :height="shape.height"
      :color="color"
      @connected="$listeners.connected"
      @connecting="$listeners.connecting"
      @unconnect="$listeners.unconnect"
    />
  </g>
</template>

<script>
import RectSelection from '../../common/rect-selection.vue';
import RectConnectors from '../../common/rect-connectors.vue';

export default {
  components: { RectConnectors, RectSelection },
  props: {
    shape: {
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
