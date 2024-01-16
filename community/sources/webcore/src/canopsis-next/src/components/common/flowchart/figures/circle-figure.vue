<template>
  <circle
    v-on="$listeners"
    :cx="centerX"
    :cy="centerY"
    :r="radius"
  />
</template>

<script>
export default {
  props: {
    diameter: {
      type: Number,
      required: true,
    },
    x: {
      type: Number,
      default: 0,
    },
    y: {
      type: Number,
      default: 0,
    },
  },
  computed: {
    radius() {
      return this.diameter / 2;
    },

    centerX() {
      return this.x + this.radius;
    },

    centerY() {
      return this.y + this.radius;
    },
  },
};
</script>
