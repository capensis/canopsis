<script>
import { LMarker } from 'vue2-leaflet';

import { GEOMAP_PANES } from '@/constants';

export default {
  extends: LMarker,
  props: {
    pane: {
      type: String,
      default: GEOMAP_PANES.markers,
    },
  },
  mounted() {
    this.mapObject.on('contextmenu', this.openContextMenu);
  },
  beforeDestroy() {
    this.mapObject.off('contextmenu', this.openContextMenu);
  },
  methods: {
    openContextMenu(event) {
      // eslint-disable-next-line no-underscore-dangle
      this.mapObject._map.fire('contextmenu', { ...event, marker: this.mapObject });
    },
  },
};
</script>
