<template lang="pug">
  geomap(ref="map", :min-zoom="minZoom")
    geomap-control-zoom(position="topleft")
    geomap-control-layers(position="topright")

    geomap-tile-layer(
      :name="$t('map.layers.openStreetMap')",
      :url="$config.OPEN_STREET_LAYER_URL",
      layer-type="base",
      no-wrap
    )

    geomap-feature-group(
      ref="pointsFeatureGroup",
      :name="$t('map.layers.points')",
      layer-type="overlay"
    )
      geomap-marker(v-for="marker in markers", :key="marker.id", :lat-lng="marker.coordinates")
        geomap-tooltip(:options="markerTooltipOptions") {{ getTooltipContent(marker.data) }}
        geomap-icon(:icon-anchor="marker.icon.anchor", :tooltip-anchor="marker.icon.tooltipAnchor")
          v-icon(
            :style="marker.icon.style",
            :size="marker.icon.size",
            color="grey darken-2"
          ) {{ marker.icon.name }}
</template>

<script>
import { getGeomapMarkerIcon } from '@/helpers/map';

import Geomap from '@/components/common/geomap/geomap.vue';
import GeomapTileLayer from '@/components/common/geomap/geomap-tile-layer.vue';
import GeomapControlZoom from '@/components/common/geomap/geomap-control-zoom.vue';
import GeomapControlLayers from '@/components/common/geomap/geomap-control-layers.vue';
import GeomapFeatureGroup from '@/components/common/geomap/geomap-feature-group.vue';
import GeomapMarker from '@/components/common/geomap/geomap-marker.vue';
import GeomapIcon from '@/components/common/geomap/geomap-icon.vue';
import GeomapTooltip from '@/components/common/geomap/geomap-tooltip.vue';

export default {
  components: {
    Geomap,
    GeomapTileLayer,
    GeomapControlZoom,
    GeomapControlLayers,
    GeomapFeatureGroup,
    GeomapMarker,
    GeomapIcon,
    GeomapTooltip,
  },
  props: {
    points: {
      type: Array,
      required: true,
    },
    minZoom: {
      type: Number,
      default: 2,
    },
    iconSize: {
      type: Number,
      default: 34,
    },
  },
  computed: {
    markerTooltipOptions() {
      return {
        offset: [0, -this.iconSize],
        direction: 'top',
      };
    },

    markers() {
      return this.points.map(point => ({
        id: point._id,
        coordinates: [point.coordinates.lat, point.coordinates.lng],
        data: point,
        icon: getGeomapMarkerIcon(point, this.iconSize),
      }));
    },
  },
  watch: {
    points() {
      this.$nextTick(this.fitMap);
    },
  },
  mounted() {
    this.$nextTick(this.fitMap);
  },
  methods: {
    fitMap() {
      const pointsBounds = this.$refs.pointsFeatureGroup.mapObject.getBounds();

      this.$refs.map.mapObject.fitBounds(pointsBounds);
    },

    getTooltipContent({ entity, map }) {
      if (entity) {
        return entity.name;
      }

      return map.name;
    },
  },
};
</script>
