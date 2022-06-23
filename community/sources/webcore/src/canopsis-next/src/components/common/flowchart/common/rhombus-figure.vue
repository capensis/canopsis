<template lang="pug">
  points-path(
    :points="rhombusPoints",
    v-on="$listeners"
  )
</template>

<script>
import PointsPath from '../common/points-path.vue';

export default {
  components: { PointsPath },
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

    rhombusPoints() {
      return [
        { x: this.x, y: this.centerY },
        { x: this.centerX, y: this.y },
        { x: this.x + this.width, y: this.centerY },
        { x: this.centerX, y: this.y + this.height },
        { x: this.x, y: this.centerY },
      ];
    },
  },
};
</script>
