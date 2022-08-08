<template lang="pug">
  div
    mermaid-point-marker(
      v-for="point in points",
      :key="point._id",
      :x="point.x",
      :y="point.y",
      :size="markerSize"
    )
</template>

<script>
import MermaidPointMarker from './partials/mermaid-point-marker.vue';

export default {
  components: { MermaidPointMarker },
  props: {
    points: {
      type: Array,
      required: true,
    },
    markerSize: {
      type: Number,
      default: 24,
    },
  },
};
</script>
