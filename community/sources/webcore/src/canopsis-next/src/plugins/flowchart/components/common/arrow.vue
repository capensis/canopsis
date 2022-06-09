<template lang="pug">
  path(:fill="fill", :d="path")
</template>

<script>
export default {
  props: {
    fill: {
      type: String,
      required: true,
    },
    width: {
      type: Number,
      default: 20,
    },
    height: {
      type: Number,
      default: 20,
    },
  },
  computed: {
    path() {
      const halfHeight = this.height / 2;
      const firstPoint = '0 0';
      const secondPoint = `0 ${this.height}`;
      const thirdPoint = `${this.width} ${halfHeight}`;

      return `M${firstPoint} ${secondPoint} ${thirdPoint} ${firstPoint} Z`;
    },
  },
};
</script>
