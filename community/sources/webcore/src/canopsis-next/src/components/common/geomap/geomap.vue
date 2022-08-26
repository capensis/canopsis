<script>
import { LMap } from 'vue2-leaflet';
import { Map, Icon } from 'leaflet';
import { GestureHandling } from 'leaflet-gesture-handling';

import locationUrl from '@/assets/images/location.svg';

Map.mergeOptions({
  attributionControl: false,
  zoomControl: false,
  gestureHandling: true,
});
Map.addInitHook('addHandler', 'gestureHandling', GestureHandling);

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
  watch: {
    disabled(value) {
      if (value) {
        this.disableInteraction();
      } else {
        this.enableInteraction();
      }
    },
  },
  mounted() {
    if (this.disabled) {
      this.disableInteraction();
    }
  },
  methods: {
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

.leaflet {
  &-pane,
  &-control,
  &-top,
  &-bottom {
    z-index: unset;
  }
}
</style>
