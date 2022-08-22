<template lang="pug">
  geomap.geomap-preview(ref="map", :min-zoom="minZoom", :disabled="!!activePoint")
    geomap-control-zoom(position="topleft", :disabled="!!activePoint")
    geomap-control-layers(position="topright", :disabled="!!activePoint")

    geomap-tile-layer(
      :name="$t('map.layers.openStreetMap')",
      :url="$config.OPEN_STREET_LAYER_URL",
      layer-type="base",
      no-wrap
    )

    geomap-cluster-group(
      ref="pointsFeatureGroup",
      :name="$t('map.layers.points')",
      layer-type="overlay"
    )
      geomap-marker(
        v-for="marker in markers",
        :key="marker.id",
        :lat-lng="marker.coordinates",
        @click="openPopup(marker.data, $event)"
      )
        geomap-icon(:icon-anchor="marker.icon.anchor", :tooltip-anchor="marker.icon.tooltipAnchor")
          point-icon(
            :style="marker.icon.style",
            :entity="marker.data.entity",
            :size="marker.icon.size",
            :color-indicator="colorIndicator",
            :pbehavior-enabled="pbehaviorEnabled",
            @show:map="$emit('show:map', $event)"
          )
    v-menu(
      v-if="activePoint",
      :value="true",
      :position-x="positionX",
      :position-y="positionY",
      :close-on-content-click="false",
      ignore-click-outside,
      offset-overflow,
      offset-x,
      absolute,
      top
    )
      point-popup(
        :point="activePoint",
        :template="popupTemplate",
        :actions="popupActions",
        @show:map="showLinkedMap",
        @close="closePopup"
      )
</template>

<script>
import { getGeomapMarkerIcon } from '@/helpers/map';

import Geomap from '@/components/common/geomap/geomap.vue';
import GeomapTileLayer from '@/components/common/geomap/geomap-tile-layer.vue';
import GeomapControlZoom from '@/components/common/geomap/geomap-control-zoom.vue';
import GeomapControlLayers from '@/components/common/geomap/geomap-control-layers.vue';
import GeomapClusterGroup from '@/components/common/geomap/geomap-cluster-group.vue';
import GeomapMarker from '@/components/common/geomap/geomap-marker.vue';
import GeomapIcon from '@/components/common/geomap/geomap-icon.vue';
import GeomapTooltip from '@/components/common/geomap/geomap-tooltip.vue';

import PointIcon from './point-icon.vue';
import PointPopup from './point-popup.vue';

export default {
  components: {
    Geomap,
    GeomapTileLayer,
    GeomapControlZoom,
    GeomapControlLayers,
    GeomapClusterGroup,
    GeomapMarker,
    GeomapIcon,
    GeomapTooltip,
    PointIcon,
    PointPopup,
  },
  props: {
    map: {
      type: Object,
      required: true,
    },
    minZoom: {
      type: Number,
      default: 2,
    },
    iconSize: {
      type: Number,
      default: 24,
    },
    popupTemplate: {
      type: String,
      required: false,
    },
    popupActions: {
      type: Boolean,
      default: false,
    },
    alarmsColumns: {
      type: Array,
      required: false,
    },
    colorIndicator: {
      type: String,
      required: false,
    },
    pbehaviorEnabled: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      positionX: 0,
      positionY: 0,
      activePoint: undefined,
    };
  },
  computed: {
    markers() {
      return this.map.parameters?.points?.map(point => ({
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
      if (this.map.parameters?.points) {
        const pointsBounds = this.$refs.pointsFeatureGroup.mapObject.getBounds();

        this.$refs.map.mapObject.fitBounds(pointsBounds);
      }
    },

    openPopup(point) {
      const { x: containerX, y: containerY } = this.$refs.map.$el.getBoundingClientRect();
      const { x, y } = this.$refs.map.mapObject.latLngToContainerPoint(point.coordinates);

      this.positionX = x + containerX;
      this.positionY = y + containerY - this.iconSize;
      this.activePoint = point;
    },

    closePopup() {
      this.activePoint = undefined;
    },

    showLinkedMap() {
      this.$emit('show:map', this.activePoint.map);
      this.closePopup();
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

<style lang="scss">
.geomap-preview {
  min-height: 700px;
}
</style>
