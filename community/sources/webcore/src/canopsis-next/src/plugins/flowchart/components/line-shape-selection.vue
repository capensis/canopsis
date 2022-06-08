<template lang="pug">
  g
    line(
      :x1="x1",
      :y1="y1",
      :x2="x2",
      :y2="y2",
      :stroke="color",
      fill="transparent",
      stroke-width="1",
      stroke-dasharray="4 4",
      pointer-events="none"
    )
    circle(
      :cx="x1",
      :cy="y1",
      :fill="color",
      :r="cornerRadius",
      cursor="pointer",
      @mousedown.stop="startMoveStart()"
    )
    circle(
      :cx="x2",
      :cy="y2",
      :fill="color",
      :r="cornerRadius",
      cursor="pointer",
      @mousedown.stop="startMoveEnd()"
    )
</template>

<script>
export default {
  props: {
    x1: {
      type: Number,
      required: true,
    },
    y1: {
      type: Number,
      required: true,
    },
    x2: {
      type: Number,
      required: true,
    },
    y2: {
      type: Number,
      required: true,
    },
    color: {
      type: String,
      default: 'blue',
    },
    cornerRadius: {
      type: Number,
      default: 4,
    },
  },
  methods: {
    startMoveStart() {
      this.$emit('start:move-from');
    },

    startMoveEnd() {
      this.$emit('start:move-to');
    },
  },
};
</script>
