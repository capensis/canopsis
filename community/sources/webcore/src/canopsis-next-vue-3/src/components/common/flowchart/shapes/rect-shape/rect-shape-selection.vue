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
      @update="$emit('update', $event)"
    />
    <rect-connectors
      v-if="connecting"
      :x="shape.x"
      :y="shape.y"
      :width="shape.width"
      :height="shape.height"
      :color="color"
      @connected="$emit('connected', $event)"
      @connecting="$emit('connecting', $event)"
      @unconnect="$emit('unconnect', $event)"
    />
  </g>
</template>

<script>
import RectConnectors from '../../common/rect-connectors.vue';
import RectSelection from '../../common/rect-selection.vue';

export default {
  components: { RectConnectors, RectSelection },
  emits: ['connected', 'connecting', 'unconnect', 'update'],
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
