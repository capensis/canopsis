<template>
  <path
    v-on="$listeners"
    :d="path"
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
      const heightWithoutRadius = this.height - (this.ellipseRadius * 2);

      return [
        /* Move to top left corner */
        `M ${this.x} ${this.y + this.ellipseRadius}`,
        /* Bottom arc to top left corner */
        `a ${this.halfWidth} ${this.ellipseRadius} 0 1, 0 ${this.width}, 0`,
        /* Top arc to top right corner */
        `a ${this.halfWidth} ${this.ellipseRadius} 0 1, 0 ${-this.width}, 0`,
        /* Line to bottom left corner */
        `l 0 ${heightWithoutRadius}`,
        /* Bottom arc to bottom right corner */
        `a ${this.halfWidth},${this.ellipseRadius} 0 1, 0 ${this.width}, 0`,
        /* Line to bottom right corner */
        `l 0 ${-heightWithoutRadius}`,
      ].join('');
    },
  },
};
</script>
