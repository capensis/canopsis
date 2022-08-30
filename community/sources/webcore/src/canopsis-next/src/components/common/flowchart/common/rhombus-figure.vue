<template lang="pug">
  path(v-on="$listeners", :d="rhombusPath")
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

    centerX() {
      return this.x + this.halfWidth;
    },

    centerY() {
      return this.y + this.halfHeight;
    },

    rhombusPath() {
      return [
        `M ${this.x} ${this.centerY}`,
        `L ${this.centerX} ${this.y}`,
        `L ${this.x + this.width} ${this.centerY}`,
        `L ${this.centerX} ${this.y + this.height}`,
        'Z',
      ].join('');
    },
  },
};
</script>
