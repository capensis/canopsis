<template lang="pug">
  path(
    :d="path",
    fill="transparent",
    pointer-events="all",
    v-on="$listeners"
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
      return this.points.map(({ x, y }, index) => `${index === 0 ? 'M' : 'L'} ${x} ${y}`).join(' ');
    },
  },
};
</script>
