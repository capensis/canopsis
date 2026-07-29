<template>
  <geomap-control :position="position">
    <c-layout column>
      <v-btn
        :disabled="disabled || zoomInDisabled"
        class="mb-1"
        color="secondary"
        icon="add"
        @click="zoomIn"
      />
      <v-btn
        :disabled="disabled || zoomOutDisabled"
        class="mb-1"
        color="secondary"
        icon="remove"
        @click="zoomOut"
      />
    </c-layout>
  </geomap-control>
</template>

<script>
import GeomapControl from './geomap-control.vue';

export default {
  inject: ['geomapRoot'],
  components: { GeomapControl },

  props: {
    position: {
      type: String,
      default: 'topleft',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },

  data() {
    return {
      zoom: 0,
    };
  },

  computed: {
    map() {
      return this.geomapRoot.mapObject;
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

  beforeUnmount() {
    this.map?.off('zoom', this.setCurrentZoom);
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
