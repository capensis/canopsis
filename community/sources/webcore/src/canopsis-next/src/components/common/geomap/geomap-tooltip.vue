<script>
import { LTooltip } from 'vue2-leaflet';
import { Tooltip } from 'leaflet';

Tooltip.mergeOptions({
  className: 'elevation-3 grey darken-2 white--text',
});

export default LTooltip;
</script>

<style lang="scss">
.leaflet-tooltip {
  border-radius: 2px;
  border: none;
  font-size: 12px;
  padding: 5px 8px;

  &:before {
    content: none;
  }
}
</style>
