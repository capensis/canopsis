<template>
  <path
    :fill="fill"
    :d="path"
  />
</template>

<script>
export default {
  props: {
    fill: {
      type: String,
      required: false,
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

      return [
        /* Move to left bottom corner */
        'M0 0',
        /* Line to bottom center */
        `L${this.width / 4} ${halfHeight}`,
        /* Line to bottom right corner */
        `L0 ${this.height}`,
        /* Line to top corner */
        `L${this.width} ${halfHeight}`,
        /* Close path */
        'Z',
      ].join('');
    },
  },
};
</script>
