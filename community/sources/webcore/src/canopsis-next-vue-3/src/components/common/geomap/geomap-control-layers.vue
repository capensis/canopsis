<template>
  <geomap-control :position="position">
    <div
@contextmenu.stop=""
@click.stop=""
@dblclick.stop=""
@mousemove.stop="">
      <v-expansion-panels>
        <v-expansion-panel class="geomap-layers-control" color="grey">
          <v-expansion-panel-title class="geomap-layers-control__title">
            <span class="v-label">
              {{ $t('geomap.layers') }}
            </span>
          </v-expansion-panel-title>
          <v-expansion-panel-text>
            <v-divider />
            <c-layout class="pa-2" column>
              <c-radio-group
                :model-value="activeLayerId"
                class="mt-0"
                color="primary"
                hide-details
                column
                @update:model-value="enableLayer"
              >
                <v-radio
                  v-for="layer in layers"
                  :key="layer.name"
                  :label="layer.name"
                  :value="getLayerId(layer.mapObject)"
                  color="primary"
                />
              </c-radio-group>
              <template v-if="overlays.length">
                <v-divider class="my-2" />
                <c-checkbox
                  v-for="overlay in overlays"
                  :key="overlay.name"
                  :model-value="isLayerActive(overlay.mapObject)"
                  :label="overlay.name"
                  class="mt-0 pt-0"
                  color="primary"
                  hide-details
                  @update:model-value="enableOverlay(overlay, $event)"
                />
              </template>
            </c-layout>
          </v-expansion-panel-text>
        </v-expansion-panel>
      </v-expansion-panels>
    </div>
  </geomap-control>
</template>

<script>
import GeomapControl from './geomap-control.vue';

export default {
  inject: ['geomapRoot'],
  components: { GeomapControl },

  props: {
    position: {
      type: String,
      default: 'topright',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },

  data() {
    return {
      allLayers: [],
      activeLayersIds: [],
    };
  },

  computed: {
    map() {
      return this.geomapRoot.mapObject;
    },

    layers() {
      return this.allLayers.filter(({ layerType, name }) => name && layerType === 'base');
    },

    overlays() {
      return this.allLayers.filter(({ layerType, name }) => name && layerType === 'overlay');
    },

    activeLayerId() {
      const activeLayer = this.layers.find(({ mapObject }) => this.isLayerActive(mapObject));

      return activeLayer ? this.getLayerId(activeLayer.mapObject) : undefined;
    },
  },

  watch: {
    allLayers: 'setActiveLayers',
  },

  mounted() {
    this.geomapRoot.registerLayerControl(this);
    this.map.on('layeradd', this.setActiveLayers, this);
    this.map.on('layerremove', this.setActiveLayers, this);
    this.setActiveLayers();
  },

  beforeUnmount() {
    this.map.off('layeradd', this.setActiveLayers, this);
    this.map.off('layerremove', this.setActiveLayers, this);
    this.geomapRoot.unregisterLayerControl(this);
  },

  methods: {
    addLayerDescriptor(layer) {
      this.allLayers.push(layer);
      this.setActiveLayers();
    },

    removeLayerDescriptor(layer) {
      this.allLayers = this.allLayers.filter(item => item !== layer);
      this.setActiveLayers();
    },

    isLayerActive(layer) {
      return this.activeLayersIds.includes(this.getLayerId(layer));
    },

    getLayerId(layer) {
      // eslint-disable-next-line no-underscore-dangle
      return layer._leaflet_id;
    },

    setActiveLayers() {
      this.activeLayersIds = this.allLayers.reduce((acc, { mapObject }) => {
        if (this.map.hasLayer(mapObject)) {
          acc.push(this.getLayerId(mapObject));
        }

        return acc;
      }, []);
    },

    addLayerToMap(layer) {
      this.map.addLayer(layer);
    },

    removeLayerFromMap(layer) {
      this.map.removeLayer(layer);
    },

    removeActiveLayers() {
      this.layers.forEach(({ mapObject }) => {
        if (this.isLayerActive(mapObject)) {
          this.removeLayerFromMap(mapObject);
        }
      });
    },

    enableLayer(layerId) {
      const layer = this.layers.find(({ mapObject }) => this.getLayerId(mapObject) === layerId);

      if (!layer) {
        return;
      }

      this.removeActiveLayers();

      this.addLayerToMap(layer.mapObject);
      this.setActiveLayers();
    },

    enableOverlay(overlay, value) {
      if (value) {
        this.addLayerToMap(overlay.mapObject);
      } else {
        this.removeLayerFromMap(overlay.mapObject);
      }

      this.setActiveLayers();
    },
  },
};
</script>

<style lang="scss">
.geomap-layers-control {
  min-width: 200px;

  &__title {
    color: white;

    .v-label {
      color: inherit;
      opacity: 1;
    }
  }
}
</style>
