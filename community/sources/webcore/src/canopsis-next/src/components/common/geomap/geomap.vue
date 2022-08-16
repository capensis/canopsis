<script>
import { LMap } from 'vue2-leaflet';
import { Map } from 'leaflet';

Map.mergeOptions({
  attributionControl: false,
});

export default LMap;
</script>

<style lang="scss">
@import "~leaflet/dist/leaflet.css";

.leaflet {
  &-pane,
  &-control,
  &-top,
  &-bottom {
    z-index: 7;
  }
}
</style>
