<template>
  <div ref="tooltipElement" style="display: none">
    <slot />
  </div>
</template>

<script>
import { Tooltip } from 'leaflet';

Tooltip.mergeOptions({
  className: 'elevation-3 bg-grey-darken-2 text-white',
});

export default {
  inject: ['geomapMarker'],

  props: {
    options: {
      type: Object,
      default: () => ({}),
    },
  },

  mounted() {
    this.geomapMarker.mapObject.bindTooltip(this.$refs.tooltipElement, this.options);
  },

  beforeUnmount() {
    this.geomapMarker.mapObject.unbindTooltip();
  },
};
</script>

<style lang="scss">
.leaflet-tooltip {
  border-radius: 2px;
  border: none;
  font-size: 12px;
  padding: 5px 8px;

  &:before {
    content: none;
  }
}
</style>
