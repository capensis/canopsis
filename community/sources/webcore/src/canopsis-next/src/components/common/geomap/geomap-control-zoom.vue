<template lang="pug">
  v-layout(column)
    v-btn.secondary.ma-0.mb-1(:disabled="zoomInDisabled", icon, dark, @click="zoomIn")
      v-icon add
    v-btn.secondary.ma-0(:disabled="zoomOutDisabled", dark, icon, @click="zoomOut")
      v-icon remove
</template>

<script>
import GeomapControl from './geomap-control.vue';

export default {
  extends: GeomapControl,
  data() {
    return {
      mapObject: undefined,
    };
  },
  computed: {
    map() {
      // eslint-disable-next-line no-underscore-dangle
      return this.mapObject?._map;
    },

    zoomInDisabled() {
      if (!this.map) {
        return false;
      }

      return this.map.getZoom() === this.map.getMaxZoom();
    },

    zoomOutDisabled() {
      if (!this.map) {
        return false;
      }

      return this.map.getZoom() === this.map.getMinZoom();
    },
  },
  methods: {
    zoomIn(event) {
      if (this.map.getZoom() < this.map.getMaxZoom()) {
        this.map.zoomIn(this.map.options.zoomDelta * (event.shiftKey ? 3 : 1));
      }
    },

    zoomOut(event) {
      if (this.map.getZoom() > this.map.getMinZoom()) {
        this.map.zoomOut(this.map.options.zoomDelta * (event.shiftKey ? 3 : 1));
      }
    },
  },
};
</script>
