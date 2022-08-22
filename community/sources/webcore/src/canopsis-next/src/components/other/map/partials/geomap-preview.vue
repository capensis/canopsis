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
        v-for="{ coordinates, id, data, icon } in markers",
        :key="id",
        :lat-lng="coordinates",
        @click="openPopup(data, $event)"
      )
        geomap-icon(:icon-anchor="icon.anchor")
          point-icon(
            :style="icon.style",
            :entity="data.entity",
            :size="icon.size",
            :color-indicator="colorIndicator",
            :pbehavior-enabled="pbehaviorEnabled"
          )
    v-menu(
      v-if="activePoint",
      :value="true",
      :position-x="positionX",
      :position-y="positionY",
      :close-on-content-click="false",
      ignore-click-outside,
      absolute,
      top
    )
      point-popup(
        :point="activePoint",
        :template="popupTemplate",
        :color-indicator="colorIndicator",
        :actions="popupActions",
        @show:alarms="showAlarms",
        @show:map="showLinkedMap",
        @close="closePopup"
      )
</template>

<script>
import { getGeomapMarkerIconOptions } from '@/helpers/map';

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
        icon: getGeomapMarkerIconOptions(point, this.iconSize),
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

    showAlarms() {
      this.$emit('show:alarms', this.activePoint);
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
