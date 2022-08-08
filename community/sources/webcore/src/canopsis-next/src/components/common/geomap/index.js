import { LMap, LTileLayer, LMarker, LControl, LIcon, LLayerGroup, LPopup } from 'vue2-leaflet';
import { Map, Icon } from 'leaflet';
import iconRetinaUrl from 'leaflet/dist/images/marker-icon-2x.png';
import iconUrl from 'leaflet/dist/images/marker-icon.png';
import shadowUrl from 'leaflet/dist/images/marker-shadow.png';

import 'leaflet/dist/leaflet.css';
import './style.scss';

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

export const Geomap = LMap;
export const GeomapTileLayer = LTileLayer;
export const GeomapMarker = LMarker;
export const GeomapControl = LControl;
export const GeomapIcon = LIcon;
export const GeomapLayerGroup = LLayerGroup;
export const GeomapPopup = LPopup;
