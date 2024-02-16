<template>
  <c-zoom-overlay>
    <v-layout
      class="geomap-editor"
      column
    >
      <geomap
        class="geomap-editor__map mb-2"
        ref="map"
        :style="mapStyles"
        :disabled="shown"
        :min-zoom="minZoom"
        :options="mapOptions"
        :center.sync="center"
        @click="openAddingPointDialogByClick"
        @dblclick="openAddingPointDialog"
      >
        <geomap-control-zoom
          position="topleft"
          :disabled="shown"
        />
        <geomap-control-layers position="topright" />
        <geomap-contextmenu
          ref="contextmenu"
          :disabled="shown"
          :items="mapItems"
          :marker-items="markerItems"
        />
        <geomap-control position="topleft">
          <v-tooltip
            right
            max-width="unset"
            min-width="max-content"
          >
            <template #activator="{ on }">
              <v-btn
                class="secondary ma-0"
                v-on="on"
                :class="{ 'lighten-4': !addOnClick }"
                :disabled="shown"
                icon
                dark
                @click="toggleAddingMode"
              >
                <v-icon>add_location</v-icon>
              </v-btn>
            </template>
            <span>{{ $t('map.toggleAddingPointMode') }}</span>
          </v-tooltip>
        </geomap-control>
        <geomap-control position="bottomright">
          <c-help-icon
            :text="$t('geomap.panzoom.helpText')"
            size="32"
            color="secondary"
            icon="help"
            top
          />
        </geomap-control>
        <geomap-tile-layer
          :name="$t('map.layers.openStreetMap')"
          :url="$config.OPEN_STREET_LAYER_URL"
          :visible="true"
          layer-type="base"
          no-wrap
        />
        <geomap-cluster-group
          ref="pointsFeatureGroup"
          :name="$t('map.layers.points')"
          layer-type="overlay"
        >
          <geomap-marker
            v-for="{ coordinates, id, data, icon } in markers"
            :key="id"
            :lat-lng="coordinates"
            :options="{ data }"
            :draggable="!data.is_entity_coordinates && !shown"
            @dragend="finishMovingMarker"
            @dblclick="openEditPointForm(data)"
          >
            <geomap-icon :icon-anchor="icon.anchor">
              <point-icon
                :style="icon.style"
                :entity="data.entity"
                :size="icon.size"
              />
            </geomap-icon>
          </geomap-marker>
        </geomap-cluster-group>
        <geomap-marker
          v-if="pointDialog"
          :lat-lng="placeholderPoint.coordinates"
        >
          <geomap-icon :icon-anchor="placeholderPoint.icon.anchor">
            <point-icon
              :style="placeholderPoint.icon.style"
              :entity="placeholderPoint.data.entity"
              :size="placeholderPoint.icon.size"
            />
          </geomap-icon>
        </geomap-marker>
        <point-form-dialog-menu
          :value="shown"
          :point="pointDialog"
          :position-x="clientX"
          :position-y="clientY"
          :editing="!!editingPoint"
          :exists-entities="existsEntities"
          coordinates
          @cancel="closePointDialog"
          @submit="submitPointDialog"
          @remove="showRemovePointModal"
          @fly:coordinates="handleUpdateCoordinates"
        />
      </geomap>
      <v-messages
        v-if="hasChildrenError"
        :value="errorMessages"
        color="error"
      />
    </v-layout>
  </c-zoom-overlay>
</template>

<script>

import { COLORS } from '@/config';
import { MODALS } from '@/constants';

import { uid } from '@/helpers/uid';
import { geomapPointToForm } from '@/helpers/entities/map/form';
import { getGeomapMarkerIconOptions } from '@/helpers/entities/map/list';

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

import PointFormDialogMenu from './point-form-dialog-menu.vue';

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
    PointFormDialogMenu,
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
      clientX: 0,
      clientY: 0,
      addingPoint: undefined,
      editingPoint: undefined,
    };
  },
  computed: {
    pointDialog() {
      return this.addingPoint || this.editingPoint;
    },

    placeholderPoint() {
      const entityPoint = { entity: uid() };

      return {
        data: entityPoint,
        coordinates: this.pointDialog.coordinates,
        icon: getGeomapMarkerIconOptions(entityPoint, this.iconSize),
      };
    },

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

    visiblePoints() {
      return this.form.points.filter(point => point._id !== this.editingPoint?._id);
    },

    markers() {
      return this.visiblePoints.map(point => ({
        id: point._id,
        coordinates: point.coordinates,
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

    existsEntities() {
      return this.form.points.map(({ entity }) => entity);
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
      const { x: containerX, y: containerY } = this.$refs.map.$el.getBoundingClientRect();
      const { x, y } = this.$refs.map.mapObject.latLngToContainerPoint(coordinates);

      this.clientX = x + containerX;
      this.clientY = y + containerY;
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

      this.flyToCoordinates(this.addingPoint.coordinates);
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

      this.flyToCoordinates(this.editingPoint.coordinates);
      this.setMenuPositionByLatLng(this.editingPoint.coordinates);
    },

    openEditPointForm(point) {
      this.editingPoint = point;
      this.shown = true;

      this.flyToCoordinates(this.editingPoint.coordinates);
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

    handleUpdateCoordinates(coordinates) {
      this.pointDialog.coordinates = coordinates;
      this.flyToCoordinates(coordinates);
    },

    flyToCoordinates(coordinates) {
      this.$refs.map.mapObject.panTo(coordinates, { animate: false });
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
