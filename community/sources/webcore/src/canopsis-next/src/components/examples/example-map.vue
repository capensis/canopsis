<template lang="pug">
  geomap(style="height: 300px")
    geomap-tile-layer(:url="$config.OPEN_STREET_LAYER_URL")
    geomap-control(position="bottomleft")
      v-btn.ma-0.primary(small) Bottom left
    geomap-control(position="topright")
      v-btn.ma-0.primary(small) Top right
    geomap-control(position="topleft")
      v-btn.ma-0.primary(small) Top left
    geomap-control(position="bottomright")
      v-btn.ma-0.primary(small) Bottom right

    geomap-layer-group(name="markers-layer")
      geomap-popup Layer layer
      geomap-marker(
        v-for="marker in markers",
        :key="marker.id",
        :lat-lng="marker.coordinate",
        :options="{ data: marker.data }",
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

import Geomap from '@/components/common/geomap/geomap.vue';
import GeomapMarker from '@/components/common/geomap/geomap-marker.vue';
import GeomapTileLayer from '@/components/common/geomap/geomap-tile-layer.vue';
import GeomapLayerGroup from '@/components/common/geomap/geomap-layer-group.vue';
import GeomapIcon from '@/components/common/geomap/geomap-icon.vue';
import GeomapControl from '@/components/common/geomap/geomap-control.vue';
import GeomapPopup from '@/components/common/geomap/geomap-popup.vue';

/**
 * TODO: Component should be removed in the end feature development
 */
export default {
  components: {
    Geomap,
    GeomapPopup,
    GeomapControl,
    GeomapIcon,
    GeomapMarker,
    GeomapLayerGroup,
    GeomapTileLayer,
  },
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
