<template lang="pug">
  geomap(style="height: 300px")
    geomap-control-zoom(position="topleft")
    geomap-control-layers(position="topright")
    geomap-contextmenu(
      ref="contextmenu",
      :items="mapItems",
      :marker-items="markerItems"
    )

    geomap-tile-layer(
      name="Open street map",
      :url="$config.OPEN_STREET_LAYER_URL",
      :visible="true",
      layer-type="base"
    )
    geomap-tile-layer(
      name="Open topo map",
      url="https://{s}.tile.opentopomap.org/{z}/{x}/{y}.png",
      :visible="false",
      layer-type="base"
    )

    geomap-layer-group(name="Markers", layer-type="overlay")
      geomap-popup Layer layer
      geomap-marker(
        v-for="marker in markers",
        :key="marker.id",
        :lat-lng="marker.coordinate",
        :options="{ data: marker.data }",
        draggable,
        @dragend="finishMovingMarker",
        @click="showInformationModal"
      )

    geomap-marker(:lat-lng="[-50, 0]")
      geomap-icon(:icon-size="[100, 38]", :icon-anchor="[50, 0]")
        v-img(:src="logo", height="20")
      geomap-popup Canopsis
</template>

<script>
import { range } from 'lodash';

import defaultLogo from '@/assets/canopsis.png';

import { MODALS } from '@/constants';

import uuid from '@/helpers/uuid';

import Geomap from '@/components/common/geomap/geomap.vue';
import GeomapMarker from '@/components/common/geomap/geomap-marker.vue';
import GeomapTileLayer from '@/components/common/geomap/geomap-tile-layer.vue';
import GeomapLayerGroup from '@/components/common/geomap/geomap-layer-group.vue';
import GeomapIcon from '@/components/common/geomap/geomap-icon.vue';
import GeomapControl from '@/components/common/geomap/geomap-control.vue';
import GeomapPopup from '@/components/common/geomap/geomap-popup.vue';
import GeomapControlZoom from '@/components/common/geomap/geomap-control-zoom.vue';
import GeomapControlLayers from '@/components/common/geomap/geomap-control-layers.vue';
import GeomapContextmenu from '@/components/common/geomap/geomap-contextmenu.vue';

/**
 * TODO: Component should be removed in the end feature development
 */
export default {
  components: {
    GeomapContextmenu,
    Geomap,
    GeomapPopup,
    GeomapControl,
    GeomapIcon,
    GeomapMarker,
    GeomapLayerGroup,
    GeomapTileLayer,
    GeomapControlZoom,
    GeomapControlLayers,
  },
  data() {
    return {
      markers: range(5).map(value => ({
        id: value,
        coordinate: [value * 10, 50 + value * 10],
        data: { id: uuid() },
      })),
    };
  },
  computed: {
    logo() {
      return defaultLogo;
    },

    mapItems() {
      return [{ text: 'Add point', action: this.addPoint }];
    },

    markerItems() {
      return [
        { text: 'Edit point', action: this.editPoint },
        { text: 'Remove point', action: this.removePoint },
      ];
    },
  },
  methods: {
    finishMovingMarker({ target }) {
      const { data } = target.options;
      const index = this.markers.findIndex(({ data: { id } }) => data.id === id);

      const { lat, lng } = target.getLatLng();

      this.markers[index].coordinate = [lat, lng];
    },

    addPoint({ latlng }) {
      const id = uuid();

      this.markers.push({
        id,
        coordinate: [latlng.lat, latlng.lng],
        data: { id },
      });

      this.$refs.contextmenu.close();
    },

    removePoint({ marker }) {
      const { data } = marker.options;

      this.markers = this.markers.filter(({ data: { id } }) => data.id !== id);

      this.$refs.contextmenu.close();
    },

    editPoint({ latlng, marker }) {
      const { data } = marker.options;

      this.$modals.show({
        name: MODALS.info,
        config: {
          title: `Edit marker: ${data.id}`,
          text: `latitude: ${latlng.lat}, longitude: ${latlng.lng}`,
        },
      });

      this.$refs.contextmenu.close();
    },

    showInformationModal(event) {
      const {
        latlng,
        target,
      } = event;

      const { data } = target.options;

      this.$modals.show({
        name: MODALS.info,
        config: {
          title: `Marker id: ${data.id}`,
          text: `latitude: ${latlng.lat}, longitude: ${latlng.lng}`,
        },
      });
    },
  },
};
</script>
