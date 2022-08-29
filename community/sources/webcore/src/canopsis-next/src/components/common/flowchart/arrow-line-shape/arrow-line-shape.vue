<template lang="pug">
  line-shape(
    v-on="$listeners",
    :shape="shape",
    :selected="selected",
    :corner-offset="cornerOffset",
    :readonly="readonly",
    :marker-end="`url(#${shape._id}-end)`",
    @update="$listeners.update",
    @edit:point="$emit('edit:point', $event)"
  )
    marker(
      :id="`${shape._id}-end`",
      refX="20",
      refY="10",
      markerWidth="60",
      markerHeight="60",
      markerUnits="userSpaceOnUse",
      orient="auto"
    )
      arrow-figure(:fill="shape.properties['stroke']")
</template>

<script>
import ArrowFigure from '../common/arrow-figure.vue';
import LineShape from '../line-shape/line-shape.vue';

export default {
  components: { LineShape, ArrowFigure },
  props: {
    shape: {
      type: Object,
      required: true,
    },
    selected: {
      type: Boolean,
      default: false,
    },
    cornerOffset: {
      type: Number,
      default: 0,
    },
    readonly: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
