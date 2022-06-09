<template lang="pug">
  path(
    :d="path",
    fill="transparent",
    pointer-events="all"
  )
</template>

<script>
export default {
  props: {
    points: {
      type: Array,
      required: true,
    },
    fill: {
      type: String,
      default: 'transparent',
    },
  },
  computed: {
    path() {
      return this.points.reduce((acc, { x, y }) => `${acc} ${x} ${y}`, 'M');
    },
  },
};
</script>
