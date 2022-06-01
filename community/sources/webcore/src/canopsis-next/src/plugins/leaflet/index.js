import { LMap, LTileLayer, LMarker, LControl, LIcon, LLayerGroup, LPopup } from 'vue2-leaflet';
import { Map, Icon } from 'leaflet';
import iconRetinaUrl from 'leaflet/dist/images/marker-icon-2x.png';
import iconUrl from 'leaflet/dist/images/marker-icon.png';
import shadowUrl from 'leaflet/dist/images/marker-shadow.png';

import 'leaflet/dist/leaflet.css';
import './styles/index.scss';

Map.mergeOptions({
  attributionControl: false,
});

// eslint-disable-next-line no-underscore-dangle
delete Icon.Default.prototype._getIconUrl;
Icon.Default.mergeOptions({
  iconRetinaUrl,
  iconUrl,
  shadowUrl,
});

export default {
  install(Vue) {
    Vue.component('c-map', LMap);
    Vue.component('c-map-tile-layer', LTileLayer);
    Vue.component('c-map-marker', LMarker);
    Vue.component('c-map-control', LControl);
    Vue.component('c-map-control', LControl);
    Vue.component('c-map-icon', LIcon);
    Vue.component('c-map-layer-group', LLayerGroup);
    Vue.component('c-map-popup', LPopup);
  },
};
