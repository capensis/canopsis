<template lang="pug">
  g(
    cursor="move",
    @mousedown="$emit('mousedown', $event)",
    @mouseup="$emit('mouseup', $event)",
    @click.stop=""
  )
    rect(
      :x="shape.x",
      :y="shape.y",
      :width="shape.width",
      :height="shape.height",
      :fill="shape.fill",
      pointer-events="all"
    )
    rect-shape-selection(v-if="selected", :shape="shape", @start:resize="onStartResize")
</template>

<script>
import RectShapeSelection from './rect-shape-selection.vue';

export default {
  components: { RectShapeSelection },
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
  },
  methods: {
    onStartResize(direction) {
      this.$emit('start:resize', direction);
    },
  },
};
</script>
