<template lang="pug">
  points-path(
    :points="parallelogramPoints",
    v-on="$listeners"
  )
</template>

<script>
import PointsPath from './points-path.vue';

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
    offset: {
      type: Number,
      default: 0,
    },
  },
  computed: {
    parallelogramPoints() {
      return [
        { x: this.x + this.offset, y: this.y },
        { x: this.x + this.width, y: this.y },
        { x: this.x + this.width - this.offset, y: this.y + this.height },
        { x: this.x, y: this.y + this.height },
        { x: this.x + this.offset, y: this.y },
      ];
    },
  },
};
</script>
