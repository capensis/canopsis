<template>
  <v-layout column>
    <v-btn
      class="secondary ma-0 mb-1"
      :disabled="disabled || zoomInDisabled"
      icon
      dark
      @click="zoomIn"
    >
      <v-icon>add</v-icon>
    </v-btn>
    <v-btn
      class="secondary ma-0 mb-1"
      :disabled="disabled || zoomOutDisabled"
      dark
      icon
      @click="zoomOut"
    >
      <v-icon>remove</v-icon>
    </v-btn>
  </v-layout>
</template>

<script>
import GeomapControl from './geomap-control.vue';

export default {
  extends: GeomapControl,
  props: {
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      zoom: 0,
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

      return this.zoom === this.map.getMaxZoom();
    },

    zoomOutDisabled() {
      if (!this.map) {
        return false;
      }

      return this.zoom === this.map.getMinZoom();
    },
  },
  mounted() {
    this.$nextTick(() => {
      this.zoom = this.map.getZoom();

      this.map.on('zoom', this.setCurrentZoom);
    });
  },
  beforeDestroy() {
    this.map.off('zoom', this.setCurrentZoom);
  },
  methods: {
    setCurrentZoom({ target }) {
      this.zoom = target.getZoom();
    },

    zoomIn(event) {
      if (this.zoom < this.map.getMaxZoom()) {
        this.map.zoomIn(this.map.options.zoomDelta * (event.shiftKey ? 3 : 1));
      }
    },

    zoomOut(event) {
      if (this.zoom > this.map.getMinZoom()) {
        this.map.zoomOut(this.map.options.zoomDelta * (event.shiftKey ? 3 : 1));
      }
    },
  },
};
</script>
