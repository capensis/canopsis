<script>
import { LTooltip } from 'vue2-leaflet';

export default LTooltip;
</script>

<style lang="scss">

.leaflet-tooltip {
  background: #616161;
  color: #fff;
  border-radius: 2px;
  border: none;
  font-size: 12px;
  display: inline-block;
  padding: 5px 8px;
  position: absolute;
  text-transform: initial;
  width: auto;
  box-shadow: 0 3px 1px -2px rgba(0,0,0,0.2),
    0px 2px 2px 0px rgba(0,0,0,0.14),
    0px 1px 5px 0px rgba(0,0,0,0.12);

  &:before {
    content: none;
  }
}
</style>
