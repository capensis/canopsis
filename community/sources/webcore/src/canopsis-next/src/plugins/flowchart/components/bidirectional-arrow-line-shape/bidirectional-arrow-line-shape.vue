<template lang="pug">
  line-shape(
    v-on="$listeners",
    v-field="shape",
    :selected="selected",
    :corner-offset="cornerOffset",
    :readonly="readonly",
    :marker-start="`url(#${shape.id}-start)`",
    :marker-end="`url(#${shape.id}-end)`"
  )
    marker(
      :id="`${shape.id}-end`",
      refX="20",
      refY="10",
      markerWidth="60",
      markerHeight="60",
      markerUnits="userSpaceOnUse",
      orient="auto"
    )
      arrow-figure(:fill="shape.style['stroke']")
    marker(
      :id="`${shape.id}-start`",
      refX="20",
      refY="10",
      markerWidth="60",
      markerHeight="60",
      markerUnits="userSpaceOnUse",
      orient="auto-start-reverse"
    )
      arrow-figure(:fill="shape.style['stroke']")
</template>

<script>
import { formBaseMixin } from '@/mixins/form';

import ArrowFigure from '../common/arrow-figure.vue';
import LineShape from '../line-shape/line-shape.vue';

export default {
  components: { LineShape, ArrowFigure },
  mixins: [formBaseMixin],
  model: {
    prop: 'shape',
    event: 'input',
  },
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
  methods: {
    move(newOffset, oldOffset) {
      const { points } = this.shape;

      this.updateModel({
        ...this.shape,

        points: points.map(({ x, y }) => ({
          x: (x - oldOffset.x) + newOffset.x,
          y: (y - oldOffset.y) + newOffset.y,
        })),
      });
    },

    onResize(line) {
      this.updateModel({
        ...this.shape,
        ...line,
      });
    },
  },
};
</script>
