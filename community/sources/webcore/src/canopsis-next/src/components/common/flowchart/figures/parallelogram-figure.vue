<template>
  <path
    v-on="$listeners"
    :d="parallelogramPath"
  />
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
      const width = this.width - this.parallelogramOffset;

      return [
        /* Move to top left */
        `M${this.x + this.parallelogramOffset} ${this.y}`,
        /* Line to top right */
        `l${this.width - this.parallelogramOffset} 0`,
        /* Line to bottom right */
        `l${-this.parallelogramOffset},${this.height}`,
        /* Line to bottom left */
        `l${-width} 0`,
        /* Close path */
        'Z',
      ].join('');
    },
  },
};
</script>
