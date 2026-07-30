<template>
  <div ref="mapElement" class="geomap">
    <slot v-if="mapObject" />
  </div>
</template>

<script>
import { debounce } from 'lodash';
import { Icon, Map } from 'leaflet';
import { GestureHandling } from 'leaflet-gesture-handling';

import { GEOMAP_PANES, GEOMAP_PANE_Z_INDEXES } from '@/constants';
import { VUETIFY_ANIMATION_DELAY } from '@/config';

import locationUrl from '@/assets/images/location.svg';

import 'leaflet/dist/leaflet.css';
import 'leaflet-gesture-handling/dist/leaflet-gesture-handling.css';

Map.mergeOptions({
  attributionControl: false,
  zoomControl: false,
  gestureHandling: true,
});
Map.addInitHook('addHandler', 'gestureHandling', GestureHandling);

const RESIZE_OBSERVER_DELAY = 100;

// eslint-disable-next-line no-underscore-dangle
delete Icon.Default.prototype._getIconUrl;
Icon.Default.mergeOptions({
  iconRetinaUrl: locationUrl,
  iconUrl: locationUrl,
  shadowUrl: false,
  iconSize: [34, 34],
  iconAnchor: [17, 31],
  popupAnchor: [1, -30],
  tooltipAnchor: [10, -15],
});

export default {
  provide() {
    return {
      geomapRoot: this,
      geomapParent: this,
    };
  },
  props: {
    disabled: {
      type: Boolean,
      required: false,
    },
    minZoom: {
      type: Number,
      required: false,
    },
    maxZoom: {
      type: Number,
      required: false,
    },
    zoom: {
      type: Number,
      default: 2,
    },
    center: {
      type: Array,
      default: () => [0, 0],
    },
    options: {
      type: Object,
      default: () => ({}),
    },
  },

  emits: ['click', 'dblclick', 'update:center'],

  data() {
    return {
      mapObject: undefined,
      layers: [],
      layerControls: [],
      invalidateSizeRaf: undefined,
      invalidateSizeTimeout: undefined,
      resizeObserver: undefined,
      debouncedInvalidateMapSize: undefined,
      isMapDestroyed: false,
    };
  },

  watch: {
    disabled(value) {
      if (value) {
        this.disableInteraction();
      } else {
        this.enableInteraction();
      }

      this.invalidateMapSize();
    },
  },

  created() {
    this.debouncedInvalidateMapSize = debounce(this.invalidateMapSize, RESIZE_OBSERVER_DELAY);
  },

  mounted() {
    this.mapObject = new Map(this.$refs.mapElement, {
      center: this.center,
      zoom: this.zoom,
      minZoom: this.minZoom,
      maxZoom: this.maxZoom,
      ...this.options,
    });

    this.mapObject.on('click', this.emitClick);
    this.mapObject.on('dblclick', this.emitDblClick);
    this.mapObject.on('moveend', this.updateCenter);

    Object.values(GEOMAP_PANES).forEach((pane) => {
      this.mapObject.createPane(pane).style.zIndex = GEOMAP_PANE_Z_INDEXES[pane];
    });

    this.observeMapResize();
    this.invalidateMapSize();

    if (this.disabled) {
      this.disableInteraction();
    }
  },

  beforeUnmount() {
    this.isMapDestroyed = true;

    this.cleanupInvalidateSize();
    this.disconnectResizeObserver();

    this.mapObject?.off('click', this.emitClick);
    this.mapObject?.off('dblclick', this.emitDblClick);
    this.mapObject?.off('moveend', this.updateCenter);
    this.mapObject?.remove();
  },

  methods: {
    emitClick(event) {
      this.$emit('click', event);
    },

    emitDblClick(event) {
      this.$emit('dblclick', event);
    },

    updateCenter() {
      const center = this.mapObject.getCenter();

      this.$emit('update:center', [center.lat, center.lng]);
    },

    isControlLayer(layer) {
      return layer.name && (layer.layerType === 'base' || layer.layerType === 'overlay');
    },

    addLayer(layer) {
      this.layers.push(layer);

      if (!layer.layerType || layer.layerType === 'base' || layer.layerType === 'overlay') {
        this.mapObject.addLayer(layer.mapObject);
      }

      if (this.isControlLayer(layer)) {
        this.layerControls.forEach((control) => control.addLayerDescriptor(layer));
      }
    },

    removeLayer(layer) {
      this.layers = this.layers.filter((item) => item !== layer);

      if (this.isControlLayer(layer)) {
        this.layerControls.forEach((control) => control.removeLayerDescriptor(layer));
      }

      if (this.mapObject.hasLayer(layer.mapObject)) {
        this.mapObject.removeLayer(layer.mapObject);
      }
    },

    registerLayerControl(control) {
      this.layerControls.push(control);
      this.layers.filter(this.isControlLayer).forEach((layer) => control.addLayerDescriptor(layer));
    },

    unregisterLayerControl(control) {
      this.layerControls = this.layerControls.filter((item) => item !== control);
    },

    cleanupInvalidateSize() {
      if (this.invalidateSizeRaf) {
        window.cancelAnimationFrame(this.invalidateSizeRaf);
      }

      if (this.invalidateSizeTimeout) {
        clearTimeout(this.invalidateSizeTimeout);
      }

      this.invalidateSizeRaf = undefined;
      this.invalidateSizeTimeout = undefined;
    },

    disconnectResizeObserver() {
      this.debouncedInvalidateMapSize?.cancel?.();

      if (this.resizeObserver) {
        this.resizeObserver.disconnect();
        this.resizeObserver = undefined;
      }
    },

    observeMapResize() {
      if (!this.$el) {
        return;
      }

      this.resizeObserver = new ResizeObserver(() => {
        this.debouncedInvalidateMapSize();
      });
      this.resizeObserver.observe(this.$el);
    },

    invalidateMapSize() {
      this.$nextTick(() => {
        if (this.isMapDestroyed || !this.mapObject) {
          return;
        }

        this.cleanupInvalidateSize();
        this.mapObject.invalidateSize();

        this.invalidateSizeRaf = window.requestAnimationFrame(() => {
          this.mapObject?.invalidateSize();
          this.invalidateSizeRaf = undefined;
        });

        this.invalidateSizeTimeout = setTimeout(() => {
          this.mapObject?.invalidateSize();
          this.invalidateSizeTimeout = undefined;
        }, VUETIFY_ANIMATION_DELAY);
      });
    },

    disableInteraction() {
      this.mapObject.scrollWheelZoom.disable();
      this.mapObject.dragging.disable();
      this.mapObject.touchZoom.disable();
      this.mapObject.boxZoom.disable();
      this.mapObject.keyboard.disable();
      this.mapObject.doubleClickZoom.disable();

      if (this.mapObject.tap) {
        this.mapObject.tap.disable();
      }
    },

    enableInteraction() {
      this.mapObject.scrollWheelZoom.enable();
      this.mapObject.dragging.enable();
      this.mapObject.touchZoom.enable();
      this.mapObject.boxZoom.enable();
      this.mapObject.keyboard.enable();

      if (this.mapObject.options?.doubleClickZoom !== false) {
        this.mapObject.doubleClickZoom.enable();
      }

      if (this.mapObject.tap) {
        this.mapObject.tap.enable();
      }
    },
  },
};
</script>

<style lang="scss">
.geomap,
.leaflet-container {
  width: 100%;
  min-height: inherit;
  height: 100%;
}

.leaflet {
  &-pane,
  &-control,
  &-top,
  &-bottom {
    z-index: unset;
  }

  &-gesture-handling-scroll-warning,
  &-gesture-handling-touch-warning {
    &::after {
      display: none;
    }
  }
}
</style>
