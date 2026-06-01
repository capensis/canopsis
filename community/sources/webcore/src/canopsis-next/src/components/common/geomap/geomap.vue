<script>
import { debounce } from 'lodash';
import { LMap } from 'vue2-leaflet';
import { Map, Icon } from 'leaflet';
import { GestureHandling } from 'leaflet-gesture-handling';

import { GEOMAP_PANES, GEOMAP_PANE_Z_INDEXES } from '@/constants';
import { VUETIFY_ANIMATION_DELAY } from '@/config';

import locationUrl from '@/assets/images/location.svg';

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
  extends: LMap,
  props: {
    disabled: {
      type: Boolean,
      required: false,
    },
  },
  data() {
    return {
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
    Object.values(GEOMAP_PANES).forEach((pane) => {
      this.mapObject.createPane(pane).style.zIndex = GEOMAP_PANE_Z_INDEXES[pane];
    });

    this.observeMapResize();
    this.invalidateMapSize();

    if (this.disabled) {
      this.disableInteraction();
    }
  },
  beforeDestroy() {
    this.isMapDestroyed = true;

    this.cleanupInvalidateSize();
    this.disconnectResizeObserver();
  },
  methods: {
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

      if (this.options?.doubleClickZoom !== false) {
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
@import "~leaflet/dist/leaflet.css";

.leaflet-container {
  width: 100%;
  min-height: inherit;
}

.leaflet {
  &-pane,
  &-control,
  &-top,
  &-bottom {
    z-index: unset;
  }
}
</style>
