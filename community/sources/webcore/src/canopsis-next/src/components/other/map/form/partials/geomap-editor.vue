<template lang="pug">
  geomap(
    ref="map",
    :style="mapStyles",
    :disabled="shown",
    :min-zoom="minZoom",
    :options="mapOptions",
    @click="openAddingPointDialogByClick",
    @dblclick="openAddingPointDialog"
  )
    geomap-control-zoom(position="topleft", :disabled="shown")
    geomap-control-layers(position="topright")

    geomap-contextmenu(
      ref="contextmenu",
      :disabled="shown",
      :items="mapItems",
      :marker-items="markerItems"
    )

    geomap-control(position="topleft")
      v-tooltip(attach, right, max-width="unset", min-width="max-content")
        template(#activator="{ on }")
          v-btn.secondary.ma-0(
            v-on="on",
            :class="{ 'lighten-4': !addOnClick }",
            :disabled="shown",
            icon,
            dark,
            @click="toggleAddingMode"
          )
            v-icon add_location
        span {{ $t('map.toggleAddingPointMode') }}

    geomap-tile-layer(
      :name="$t('map.layers.openStreetMap')",
      :url="$config.OPEN_STREET_LAYER_URL",
      :visible="true",
      layer-type="base",
      no-wrap
    )

    geomap-feature-group(
      ref="pointsFeatureGroup",
      :name="$t('map.layers.points')",
      layer-type="overlay"
    )
      geomap-marker(
        v-for="marker in markers",
        :key="marker.id",
        :lat-lng="marker.coordinates",
        :options="{ data: marker.data }",
        :draggable="!shown",
        @dragend="finishMovingMarker",
        @click=""
      )
        geomap-icon(:icon-anchor="marker.icon.anchor")
          v-icon(
            :style="marker.icon.style",
            :size="marker.icon.size",
            color="grey darken-2"
          ) {{ marker.icon.name }}

    v-menu(
      v-model="shown",
      :position-x="pageX",
      :position-y="pageY",
      :close-on-content-click="false",
      ignore-click-outside,
      offset-overflow,
      offset-x,
      absolute
    )
      point-form-dialog(
        v-if="addingPoint || editingPoint",
        :point="addingPoint || editingPoint",
        :editing="!!editingPoint",
        coordinates,
        @cancel="closePointDialog",
        @submit="submitPointDialog",
        @remove="showRemovePointModal"
      )
</template>

<script>
import { MODALS } from '@/constants';

import { geomapPointToForm } from '@/helpers/forms/map';

import { formMixin } from '@/mixins/form';

import Geomap from '@/components/common/geomap/geomap.vue';
import GeomapTileLayer from '@/components/common/geomap/geomap-tile-layer.vue';
import GeomapControl from '@/components/common/geomap/geomap-control.vue';
import GeomapControlZoom from '@/components/common/geomap/geomap-control-zoom.vue';
import GeomapControlLayers from '@/components/common/geomap/geomap-control-layers.vue';
import GeomapContextmenu from '@/components/common/geomap/geomap-contextmenu.vue';
import GeomapFeatureGroup from '@/components/common/geomap/geomap-feature-group.vue';
import GeomapMarker from '@/components/common/geomap/geomap-marker.vue';
import GeomapIcon from '@/components/common/geomap/geomap-icon.vue';

import PointFormDialog from './point-form-dialog.vue';

export default {
  components: {
    Geomap,
    GeomapTileLayer,
    GeomapContextmenu,
    GeomapControlZoom,
    GeomapControlLayers,
    GeomapControl,
    GeomapFeatureGroup,
    GeomapMarker,
    PointFormDialog,
    GeomapIcon,
  },
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
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
  data() {
    return {
      center: undefined,
      addOnClick: false,

      shown: false,
      pageX: 0,
      pageY: 0,
      addingPoint: undefined,
      editingPoint: undefined,
    };
  },
  computed: {
    mapOptions() {
      return {
        doubleClickZoom: false,
      };
    },

    mapStyles() {
      return {
        cursor: this.addOnClick ? 'crosshair' : '',
      };
    },

    markers() {
      const halfIconSize = this.iconSize / 2;
      const pixelSize = `${this.iconSize}px`;

      return this.form.points.map(point => ({
        id: point._id,
        coordinates: [point.coordinates.lat, point.coordinates.lng],
        data: point,
        icon: {
          name: point.entity ? 'location_on' : 'link',
          style: {
            width: pixelSize,
            height: pixelSize,
            maxWidth: 'unset',
            maxHeight: 'unset',
          },
          size: this.iconSize,
          anchor: point.entity
            ? [halfIconSize, this.iconSize]
            : [halfIconSize, halfIconSize],
        },
      }));
    },

    mapItems() {
      return [{ text: this.$t('map.addPoint'), action: this.addPoint }];
    },

    markerItems() {
      return [
        { text: this.$t('map.editPoint'), action: this.editPoint },
        { text: this.$t('map.removePoint'), action: this.removePoint },
      ];
    },
  },
  mounted() {
    this.$nextTick(this.fitMap);
  },
  methods: {
    fitMap() {
      const { pointsFeatureGroup, map } = this.$refs;

      if (this.form.points.length) {
        const pointsBounds = pointsFeatureGroup.mapObject.getBounds();

        map.mapObject.fitBounds(pointsBounds);
      } else {
        map.mapObject.fitWorld();
      }
    },

    updatePointInModel(data) {
      this.updateField('points', this.form.points.map(point => (point._id === data._id ? data : point)));
    },

    addPointToModel(data) {
      this.updateField('points', [...this.form.points, data]);
    },

    removePointFromModel(data) {
      this.updateField(
        'points',
        this.form.points.filter(point => data._id !== point._id),
      );
    },

    setMenuPositionByLatLng(coordinates) {
      if (coordinates) {
        const { x: containerX, y: containerY } = this.$refs.map.$el.getBoundingClientRect();
        const { x, y } = this.$refs.map.mapObject.latLngToContainerPoint(coordinates);

        this.pageX = x + containerX;
        this.pageY = y + containerY;
      }
    },

    closeContextMenu() {
      this.$refs.contextmenu.close();
    },

    closePointDialog() {
      this.shown = false;
      this.addingPoint = undefined;
      this.editingPoint = undefined;
    },

    toggleAddingMode() {
      this.addOnClick = !this.addOnClick;
    },

    openAddingPointDialogByClick(event) {
      if (!this.addOnClick) {
        return;
      }

      this.openAddingPointDialog(event);
    },

    openAddingPointDialog({ latlng }) {
      if (this.shown) {
        return;
      }

      this.shown = true;

      this.addingPoint = geomapPointToForm({
        coordinates: {
          lat: latlng.lat,
          lng: latlng.lng,
        },
      });

      this.setMenuPositionByLatLng(this.addingPoint.coordinates);
    },

    addPoint({ latlng }) {
      this.openAddingPointDialog({ latlng });

      this.closeContextMenu();
    },

    editPoint({ marker }) {
      this.editingPoint = marker.options.data;

      this.shown = true;

      this.closeContextMenu();

      this.setMenuPositionByLatLng(this.editingPoint.coordinates);
    },

    submitPointDialog(data) {
      if (this.editingPoint) {
        this.updatePointInModel(data);
      } else {
        this.addPointToModel(data);
      }

      this.closePointDialog();
    },

    finishMovingMarker({ target }) {
      const { lat, lng } = target.getLatLng();

      this.updatePointInModel({
        ...target.options.data,
        coordinates: { lat, lng },
      });
    },

    removePoint({ marker }) {
      this.showRemovePointModal(marker.options.data);

      this.closeContextMenu();
    },

    showRemovePointModal(point) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => {
            this.removePointFromModel(point || this.editingPoint);

            this.closePointDialog();
          },
        },
      });
    },
  },
};
</script>
