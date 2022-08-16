<script>
import { LIcon } from 'vue2-leaflet';
import { Icon } from 'leaflet';

import locationUrl from '@/assets/images/location.svg';

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

export default LIcon;
</script>
