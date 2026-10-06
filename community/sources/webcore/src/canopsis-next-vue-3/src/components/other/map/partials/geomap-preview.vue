<template>
  <c-zoom-overlay>
    <geomap
      ref="map"
      :min-zoom="minZoom"
      :disabled="!!activePoint"
      class="geomap-preview"
    >
      <geomap-control-zoom :disabled="!!activePoint" position="topleft" />
      <geomap-control-layers :disabled="!!activePoint" position="topright" />
      <geomap-tile-layer
        :name="$t('map.layers.openStreetMap')"
        :url="openStreetLayerUrl"
        layer-type="base"
        no-wrap
      />
      <geomap-control position="bottomright">
        <c-help-icon
          :text="$t('geomap.panzoom.helpText')"
          size="32"
          color="secondary"
          icon="help"
          location="top"
        />
      </geomap-control>
      <geomap-cluster-group
        v-for="{ markers: layerMarkers, name, style } in layers"
        :key="name"
        :name="name"
        :cluster-style="style"
        layer-type="overlay"
      >
        <geomap-marker
          v-for="{ coordinates, id, data, icon } in layerMarkers"
          :key="id"
          :lat-lng="coordinates"
          @click="openMarkerPopup(data, $event)"
        >
          <geomap-icon :icon-anchor="icon.anchor">
            <point-icon
              :style="icon.style"
              :entity="data.entity"
              :size="icon.size"
              :color-indicator="colorIndicator"
              :pbehavior-enabled="pbehaviorEnabled"
              :pause-enabled="pauseEnabled"
            />
          </geomap-icon>
        </geomap-marker>
      </geomap-cluster-group>
      <point-popup-dialog
        v-if="activePoint"
        :point="activePoint"
        :position-x="positionX"
        :position-y="positionY"
        :popup-template="popupTemplate"
        :color-indicator="colorIndicator"
        :popup-actions="popupActions"
        @show:alarms="showAlarms"
        @show:map="showLinkedMap"
        @close="closePopup"
      />
    </geomap>
  </c-zoom-overlay>
</template>

<script>
import { LatLng, LatLngBounds } from 'leaflet';
import { groupBy } from 'lodash';

import { OPEN_STREET_LAYER_URL } from '@/config';
import { COLOR_INDICATOR_TYPES, PBEHAVIOR_TYPE_TYPES } from '@/constants';

import { getEntityColor } from '@/helpers/entities/entity/color';
import { getGeomapMarkerIconOptions } from '@/helpers/entities/map/list';

import { mapInformationPopupMixin } from '@/mixins/map/map-information-popup-mixin';

import GeomapClusterGroup from '@/components/common/geomap/geomap-cluster-group.vue';
import GeomapControlLayers from '@/components/common/geomap/geomap-control-layers.vue';
import GeomapControlZoom from '@/components/common/geomap/geomap-control-zoom.vue';
import GeomapControl from '@/components/common/geomap/geomap-control.vue';
import GeomapIcon from '@/components/common/geomap/geomap-icon.vue';
import GeomapMarker from '@/components/common/geomap/geomap-marker.vue';
import GeomapTileLayer from '@/components/common/geomap/geomap-tile-layer.vue';
import Geomap from '@/components/common/geomap/geomap.vue';

import PointIcon from './point-icon.vue';
import PointPopupDialog from './point-popup-dialog.vue';

export default {
  components: {
    PointPopupDialog,
    Geomap,
    GeomapTileLayer,
    GeomapControlZoom,
    GeomapControlLayers,
    GeomapClusterGroup,
    GeomapMarker,
    GeomapIcon,
    GeomapControl,
    PointIcon,
  },

  mixins: [mapInformationPopupMixin],

  props: {
    map: {
      type: Object,
      required: true,
    },
    minZoom: {
      type: Number,
      default: 2,
    },
  },

  computed: {
    openStreetLayerUrl() {
      return OPEN_STREET_LAYER_URL;
    },

    isStateColorIndicator() {
      return this.colorIndicator === COLOR_INDICATOR_TYPES.state;
    },

    markers() {
      return this.map.parameters?.points?.map(point => ({
        id: point._id,
        coordinates: [point.coordinates.lat, point.coordinates.lng],
        data: point,
        icon: getGeomapMarkerIconOptions(point, this.iconSize),
      }));
    },

    groupedMarkers() {
      return groupBy(this.markers ?? [], ({ data }) => {
        const { entity } = data;

        if (!entity) {
          return 'map';
        }

        const hasNonPausePbehavior = entity.pbehavior_info &&
          entity.pbehavior_info.canonical_type !== PBEHAVIOR_TYPE_TYPES.pause;

        if (this.pbehaviorEnabled && hasNonPausePbehavior) {
          return 'pbehavior';
        }

        if (this.pauseEnabled && entity.pause) {
          return 'pause';
        }

        if (this.pbehaviorEnabled && entity.pbehavior_info) {
          return 'pbehavior';
        }

        if (this.isStateColorIndicator) {
          return entity.state;
        }

        return entity.impact_state;
      });
    },

    layers() {
      const layers = {};

      if (!this.colorIndicator && !this.pbehaviorEnabled && !this.pauseEnabled) {
        return {
          points: {
            name: this.$t('map.layers.points'),
            markers: this.markers,
          },
        };
      }

      const {
        map: mapMarkers,
        pbehavior: pbehaviorMarkers,
        pause: pauseMarkers,
        ...restGroups
      } = this.groupedMarkers;

      if (mapMarkers) {
        layers.map = {
          name: this.$t('common.map', mapMarkers.length),
          markers: mapMarkers,
        };
      }

      if (pbehaviorMarkers) {
        layers.pbehavior = {
          name: this.$t('common.pbehavior'),
          markers: pbehaviorMarkers,
        };
      }

      if (pauseMarkers) {
        layers.pause = {
          name: this.$t('common.pause'),
          markers: pauseMarkers,
        };
      }

      Object.entries(restGroups).forEach(([key, markers]) => {
        const [firstMarker] = markers;

        layers[key] = {
          style: { backgroundColor: getEntityColor(firstMarker.data.entity, this.colorIndicator) },
          name: this.isStateColorIndicator
            ? this.$t(`common.stateTypes.${key}`)
            : `${this.$t('common.impactLevel')}: ${key}`,
          markers,
        };
      });

      return layers;
    },
  },

  watch: {
    'map.parameters.points': {
      handler() {
        this.$nextTick(this.fitMap);
      },
    },
  },

  mounted() {
    this.$nextTick(this.fitMap);
  },

  methods: {
    fitMap() {
      if (this.map.parameters?.points?.length) {
        const bounds = new LatLngBounds();

        this.map.parameters.points.forEach(({ coordinates }) => {
          bounds.extend(new LatLng(coordinates.lat, coordinates.lng));
        });

        this.$refs.map.mapObject.fitBounds(bounds);
      }
    },

    openMarkerPopup(point, { originalEvent }) {
      this.openPopup(point, originalEvent);
    },
  },
};
</script>

<style lang="scss">
.geomap-preview {
  height: 700px;
  min-height: 700px;
  border-radius: 5px;
  overflow: hidden;
}
</style>
