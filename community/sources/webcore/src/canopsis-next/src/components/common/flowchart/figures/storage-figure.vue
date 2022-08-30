<template lang="pug">
  path(v-on="$listeners", :d="path")
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
    radius: {
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

    ellipseRadius() {
      return this.radius * 2 <= this.height ? this.radius : this.height / 2;
    },

    path() {
      return [
        `M ${this.x},${this.y + this.ellipseRadius}`,
        `a ${this.halfWidth},${this.ellipseRadius} 0 1, 0 ${this.width}, 0`,
        `a ${this.halfWidth},${this.ellipseRadius} 0 1, 0 ${-this.width}, 0`,
        `L ${this.x},${this.y + this.height - this.ellipseRadius}`,
        `a ${this.halfWidth},${this.ellipseRadius} 0 1, 0 ${this.width}, 0`,
        `L ${this.x + this.width},${this.y + this.ellipseRadius}`,
      ].join('');
    },
  },
};
</script>
