<template lang="pug">
  path(v-on="$listeners", :d="parallelogramPath")
</template>

<script>
export default {
  props: {
    width: {
      type: Number,
      required: true,
    },
    height: {
      type: Number,
      required: true,
    },
    x: {
      type: Number,
      required: true,
    },
    y: {
      type: Number,
      required: true,
    },
    offset: {
      type: Number,
      default: 0,
    },
  },
  computed: {
    parallelogramOffset() {
      return this.offset < this.width ? this.offset : this.width;
    },

    parallelogramPath() {
      return [
        `M${this.x + this.parallelogramOffset},${this.y}`,
        `L${this.x + this.width},${this.y}`,
        `L${this.x + this.width - this.parallelogramOffset},${this.y + this.height}`,
        `L${this.x},${this.y + this.height}`,
        `L${this.x + this.parallelogramOffset},${this.y}`,
        'Z',
      ].join(' ');
    },
  },
};
</script>
