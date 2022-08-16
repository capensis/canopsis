<template lang="pug">
  div
    mermaid-point-marker(
      v-for="point in points",
      :key="point._id",
      :x="point.x",
      :y="point.y",
      :entity="point.entity",
      :size="markerSize",
      @mouseenter="onMouseEnter(point, $event)",
      @mouseleave="onMouseLeave"
    )
    v-tooltip(
      v-if="tooltipContent",
      :value="true",
      :position-x="positionX",
      :position-y="positionY",
      top
    )
      span {{ tooltipContent }}
</template>

<script>
import MermaidPointMarker from './mermaid-point-marker.vue';

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
  data() {
    return {
      positionX: 0,
      positionY: 0,
      tooltipContent: undefined,
    };
  },
  methods: {
    onMouseEnter(point, event) {
      const { entity, map } = point;

      const { top, left, width } = event.target.getBoundingClientRect();

      this.positionY = top;
      this.positionX = left + width / 2;
      this.tooltipContent = entity ? entity.name : map.name;
    },

    onMouseLeave() {
      this.tooltipContent = undefined;
    },
  },
};
</script>
