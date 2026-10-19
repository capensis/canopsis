<template>
  <path
    :d="path"
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
    offset: {
      type: Number,
      default: 0,
    },
  },
  computed: {
    halfOffset() {
      return this.offset / 2;
    },

    halfWidth() {
      return this.width / 2;
    },

    quarterWidth() {
      return this.halfWidth / 2;
    },

    path() {
      return [
        /* Move to top left */
        `M ${this.x} ${this.y}`,
        /* Move to top right */
        `l ${this.width} 0`,
        /* Move to bottom right */
        `l 0 ${this.height - this.halfOffset}`,
        /* Curve to bottom center */
        `q ${-this.quarterWidth} ${-this.offset} ${-this.halfWidth} 0`,
        /* Curve from bottom center to bottom left */
        `q ${-this.quarterWidth},${this.offset} ${-this.halfWidth} 0`,
        /* Close path */
        'Z',
      ].join('');
    },
  },
};
</script>
