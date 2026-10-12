<template>
  <line-shape
    :shape="shape"
    :readonly="readonly"
    :marker-end="`url(#${shape._id}-end)`"
    v-on="$listeners"
  >
    <marker
      :id="`${shape._id}-end`"
      refX="20"
      refY="10"
      markerWidth="60"
      markerHeight="60"
      markerUnits="userSpaceOnUse"
      orient="auto"
    >
      <arrow-figure :fill="shape.properties['stroke']" />
    </marker>
  </line-shape>
</template>

<script>
import ArrowFigure from '../../figures/arrow-figure.vue';
import LineShape from '../line-shape/line-shape.vue';

export default {
  components: { LineShape, ArrowFigure },
  props: {
    shape: {
      type: Object,
      required: true,
    },
    readonly: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
