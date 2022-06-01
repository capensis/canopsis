<template lang="pug">
  c-map(style="height: 300px")
    c-map-tile-layer(:url="$constants.MAP_LAYERS.openStreetMap")
    c-map-control(position="bottomleft")
      v-btn.ma-0.primary(small) Bottom left
    c-map-control(position="topright")
      v-btn.ma-0.primary(small) Top right
    c-map-control(position="topleft")
      v-btn.ma-0.primary(small) Top left
    c-map-control(position="bottomright")
      v-btn.ma-0.primary(small) Bottom right

    c-map-layer-group(name="markers-layer")
      c-map-popup Layer layer
      c-map-marker(
        v-for="marker in markers",
        :key="marker.id",
        :lat-lng="marker.coordinate",
        :options="{ data: marker.data }",
        @click="showInformationModal"
      )

    c-map-marker(:lat-lng="[-50, 0]")
      c-map-icon(:icon-size="[100, 38]", :icon-anchor="[50, 0]")
        v-img(:src="logo", height="20")
      c-map-popup Canopsis
</template>

<script>
import { range } from 'lodash';

import defaultLogo from '@/assets/canopsis.png';

import { MODALS } from '@/constants';

export default {
  computed: {
    logo() {
      return defaultLogo;
    },

    markers() {
      return range(5).map(value => ({
        id: value,
        coordinate: [0 + value * 10, 50 + value * 10],
        data: { id: value },
      }));
    },
  },
  methods: {
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
