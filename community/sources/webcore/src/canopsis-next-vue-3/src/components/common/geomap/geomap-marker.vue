<template>
  <div style="display: none">
    <slot />
  </div>
</template>

<script>
import { marker } from 'leaflet';

import { GEOMAP_PANES } from '@/constants';

export default {
  inject: ['geomapRoot', 'geomapParent'],

  provide() {
    return {
      geomapMarker: this,
    };
  },

  props: {
    latLng: {
      type: Array,
      required: true,
    },
    pane: {
      type: String,
      default: GEOMAP_PANES.markers,
    },
  },

  emits: ['click'],

  data() {
    return {
      mapObject: undefined,
    };
  },

  watch: {
    latLng(value) {
      this.mapObject?.setLatLng(value);
    },
  },

  beforeMount() {
    this.mapObject = marker(this.latLng, {
      pane: this.pane,
    });
  },

  mounted() {
    this.mapObject.on('click', (event) => this.$emit('click', event));
    this.mapObject.on('contextmenu', this.openContextMenu);
    this.geomapParent.addLayer(this);
  },

  beforeUnmount() {
    if (this.mapObject) {
      this.mapObject.off('contextmenu', this.openContextMenu);
      this.geomapParent.removeLayer(this);
    }
  },

  methods: {
    setIcon(icon) {
      this.mapObject?.setIcon(icon);
    },

    openContextMenu(event) {
      this.geomapRoot.mapObject.fire('contextmenu', { ...event, marker: this.mapObject });
    },
  },
};
</script>
