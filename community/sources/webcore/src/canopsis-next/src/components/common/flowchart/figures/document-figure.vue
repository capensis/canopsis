<template lang="pug">
  path(:d="path", v-on="$listeners")
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
        `M ${this.x},${this.y}`,
        `l ${this.width},0`,
        `l 0,${this.height - this.halfOffset}`,
        `q ${-this.quarterWidth},${-this.offset}  ${-this.halfWidth},0`,
        `q ${-this.quarterWidth},${this.offset}  ${-this.halfWidth},0`,
        'Z',
      ].join('');
    },
  },
};
</script>
