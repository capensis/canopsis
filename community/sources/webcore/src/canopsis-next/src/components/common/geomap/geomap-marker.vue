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
    draggable: {
      type: Boolean,
      default: false,
    },
    options: {
      type: Object,
      default: () => ({}),
    },
  },

  emits: ['click', 'dblclick', 'dragend'],

  data() {
    return {
      mapObject: undefined,
    };
  },

  watch: {
    latLng(value) {
      this.mapObject?.setLatLng(value);
    },

    draggable(value) {
      if (value) {
        this.mapObject?.dragging.enable();
      } else {
        this.mapObject?.dragging.disable();
      }
    },
  },

  beforeMount() {
    this.mapObject = marker(this.latLng, {
      pane: this.pane,
      draggable: this.draggable,
      ...this.options,
    });
  },

  mounted() {
    this.mapObject.on('click', this.emitClick);
    this.mapObject.on('dblclick', this.emitDblClick);
    this.mapObject.on('dragend', this.emitDragEnd);
    this.mapObject.on('contextmenu', this.openContextMenu);
    this.geomapParent.addLayer(this);
  },

  beforeUnmount() {
    if (this.mapObject) {
      this.mapObject.off('click', this.emitClick);
      this.mapObject.off('dblclick', this.emitDblClick);
      this.mapObject.off('dragend', this.emitDragEnd);
      this.mapObject.off('contextmenu', this.openContextMenu);
      this.geomapParent.removeLayer(this);
    }
  },

  methods: {
    emitClick(event) {
      this.$emit('click', event);
    },

    emitDblClick(event) {
      this.$emit('dblclick', event);
    },

    emitDragEnd(event) {
      this.$emit('dragend', event);
    },

    setIcon(icon) {
      this.mapObject?.setIcon(icon);
    },

    openContextMenu(event) {
      this.geomapRoot.mapObject.fire('contextmenu', { ...event, marker: this.mapObject });
    },
  },
};
</script>
