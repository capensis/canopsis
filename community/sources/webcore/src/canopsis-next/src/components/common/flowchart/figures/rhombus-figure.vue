<template>
  <path
    :d="rhombusPath"
    v-on="$listeners"
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
  },
  computed: {
    halfWidth() {
      return this.width / 2;
    },

    halfHeight() {
      return this.height / 2;
    },

    rhombusPath() {
      return [
        /* Move to left corner */
        `M${this.x} ${this.y + this.halfHeight}`,
        /* Line to top corner */
        `l${this.halfWidth} ${-this.halfHeight}`,
        /* Line to right corner */
        `l${this.halfWidth} ${this.halfHeight}`,
        /* Line to bottom corner */
        `l${-this.halfWidth} ${this.halfHeight}`,
        /* Close path */
        'Z',
      ].join('');
    },
  },
};
</script>
