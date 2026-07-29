<template>
  <div ref="iconElement" style="display: none">
    <slot />
  </div>
</template>

<script>
import { divIcon } from 'leaflet';

export default {
  inject: ['geomapMarker'],

  props: {
    iconAnchor: {
      type: Array,
      required: false,
    },
  },

  mounted() {
    this.updateIcon();
  },

  updated() {
    this.updateIcon();
  },

  methods: {
    updateIcon() {
      this.geomapMarker.setIcon(
        divIcon({
          html: this.$refs.iconElement.innerHTML,
          className: '',
          iconAnchor: this.iconAnchor,
        }),
      );
    },
  },
};
</script>
