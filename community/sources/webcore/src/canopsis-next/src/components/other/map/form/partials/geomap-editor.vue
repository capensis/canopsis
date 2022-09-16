<template lang="pug">
  v-layout.geomap-editor(column)
    geomap.geomap-editor__map.mb-2(
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
        v-tooltip(right, max-width="unset", min-width="max-content")
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

      geomap-control(position="bottomright")
        c-help-icon(size="32", color="secondary", icon="help", top)
          div.pre-wrap(v-html="$t('geomap.panzoom.helpText')")

      geomap-tile-layer(
        :name="$t('map.layers.openStreetMap')",
        :url="$config.OPEN_STREET_LAYER_URL",
        :visible="true",
        layer-type="base",
        no-wrap
      )

      geomap-cluster-group(
        ref="pointsFeatureGroup",
        :name="$t('map.layers.points')",
        :disable-clustering-at-zoom="maxClusteringZoom",
        layer-type="overlay"
      )
        geomap-marker(
          v-for="{ coordinates, id, data, icon } in markers",
          :key="id",
          :lat-lng="coordinates",
          :options="{ data }",
          :draggable="!shown",
          @dragend="finishMovingMarker",
          @click=""
        )
          geomap-icon(:icon-anchor="icon.anchor")
            point-icon(:style="icon.style", :entity="data.entity", :size="icon.size")

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
    v-messages(v-if="hasChildrenError", :value="errorMessages", color="error")
</template>

<script>
import { COLORS } from '@/config';

import { MODALS } from '@/constants';

import { geomapPointToForm } from '@/helpers/forms/map';
import { getGeomapMarkerIconOptions } from '@/helpers/map';

import { formMixin, validationChildrenMixin } from '@/mixins/form';

import Geomap from '@/components/common/geomap/geomap.vue';
import GeomapTileLayer from '@/components/common/geomap/geomap-tile-layer.vue';
import GeomapControl from '@/components/common/geomap/geomap-control.vue';
import GeomapControlZoom from '@/components/common/geomap/geomap-control-zoom.vue';
import GeomapControlLayers from '@/components/common/geomap/geomap-control-layers.vue';
import GeomapContextmenu from '@/components/common/geomap/geomap-contextmenu.vue';
import GeomapClusterGroup from '@/components/common/geomap/geomap-cluster-group.vue';
import GeomapMarker from '@/components/common/geomap/geomap-marker.vue';
import GeomapIcon from '@/components/common/geomap/geomap-icon.vue';
import PointIcon from '@/components/other/map/partials/point-icon.vue';

import PointFormDialog from './point-form-dialog.vue';

export default {
  inject: ['$validator'],
  components: {
    Geomap,
    GeomapTileLayer,
    GeomapContextmenu,
    GeomapControlZoom,
    GeomapControlLayers,
    GeomapControl,
    GeomapClusterGroup,
    GeomapMarker,
    PointFormDialog,
    GeomapIcon,
    PointIcon,
  },
  mixins: [formMixin, validationChildrenMixin],
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
    maxClusteringZoom: {
      type: Number,
      default: 12,
    },
    iconSize: {
      type: Number,
      default: 24,
    },
    name: {
      type: String,
      default: 'parameters',
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
    errorMessages() {
      return [this.$t('geomap.errors.pointsRequired')];
    },

    mapOptions() {
      return {
        doubleClickZoom: false,
      };
    },

    mapStyles() {
      return {
        cursor: this.addOnClick ? 'crosshair' : '',
        borderColor: this.hasChildrenError ? COLORS.error : undefined,
      };
    },

    markers() {
      return this.form.points.map(point => ({
        id: point._id,
        coordinates: [point.coordinates.lat, point.coordinates.lng],
        data: point,
        icon: getGeomapMarkerIconOptions(point, this.iconSize),
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

    pointsFieldName() {
      return `${this.name}.points`;
    },
  },
  watch: {
    form: {
      deep: true,
      handler() {
        if (this.hasChildrenError) {
          this.$validator.validate(this.pointsFieldName);
        }
      },
    },
  },
  beforeDestroy() {
    this.detachRules();
  },
  mounted() {
    this.attachRequiredRule();
    this.$nextTick(this.fitMap);
  },
  methods: {
    attachRequiredRule() {
      this.$validator.attach({
        name: this.pointsFieldName,
        rules: 'required:true',
        getter: () => !!this.form.points.length,
        vm: this,
      });
    },

    detachRules() {
      this.$validator.detach(this.pointsFieldName);
    },

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

<style lang="scss">
$borderColor: #e5e5e5;

.geomap-editor {
  &__map {
    min-height: 500px;
    border: 1px solid $borderColor;
  }
}
</style>
