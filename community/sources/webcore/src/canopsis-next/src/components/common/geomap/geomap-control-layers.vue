<template>
  <div
    @contextmenu.stop=""
    @click.stop=""
    @dblclick.stop=""
    @mousemove.stop=""
  >
    <v-expansion-panels>
      <v-expansion-panel
        class="geomap-layers-control"
        color="grey"
      >
        <v-expansion-panel-header>
          <span class="v-label">{{ $t('geomap.layers') }}</span>
        </v-expansion-panel-header>
        <v-expansion-panel-content>
          <v-divider />
          <v-layout
            class="pa-2"
            column
          >
            <v-radio-group
              class="mt-0"
              :value="activeLayer"
              color="primary"
              column
              hide-details
              @change="enableLayer"
            >
              <v-radio
                v-for="layer in layers"
                :key="layer.name"
                :label="layer.name"
                :value="layer"
                color="primary"
              />
            </v-radio-group>
            <template v-if="overlays.length">
              <v-divider class="my-2" />
              <v-checkbox
                class="mt-0 pt-0"
                v-for="overlay in overlays"
                :key="overlay.name"
                :input-value="isLayerActive(overlay.layer)"
                :label="overlay.name"
                hide-details
                @change="enableOverlay(overlay, $event)"
              />
            </template>
          </v-layout>
        </v-expansion-panel-content>
      </v-expansion-panel>
    </v-expansion-panels>
  </div>
</template>

<script>
import { LControlLayers, optionsMerger, propsBinder } from 'vue2-leaflet';
import { Control, control, DomUtil } from 'leaflet';

Control.Layers = Control.Layers.extend({
  element: undefined,
  setElement(el) {
    this.element = el;
  },

  _initLayout() {
    // eslint-disable-next-line no-underscore-dangle
    this._container = DomUtil.create('div', 'leaflet-control');

    if (this.element) {
      // eslint-disable-next-line no-underscore-dangle
      this._container.appendChild(this.element);
    }
  },
  _update() {},
  expand() {},
  collapse() {},
});
control.layers = (...args) => new Control.Layers(...args);

export default {
  ...LControlLayers,

  data() {
    return {
      mapObject: undefined,
      allLayers: [],
      activeLayersIds: [],
    };
  },
  computed: {
    map() {
      // eslint-disable-next-line no-underscore-dangle
      return this.mapObject._map;
    },

    layers() {
      return this.allLayers.filter(({ overlay }) => !overlay);
    },

    overlays() {
      return this.allLayers.filter(({ overlay }) => overlay);
    },

    activeLayer() {
      return this.layers.find(({ layer }) => this.isLayerActive(layer));
    },
  },
  watch: {
    layers: 'setActiveLayers',
    overlays: 'setActiveLayers',
  },
  mounted() {
    const options = optionsMerger(
      {
        ...this.controlOptions,
        collapsed: this.collapsed,
        autoZIndex: this.autoZIndex,
        hideSingleBase: this.hideSingleBase,
        sortLayers: this.sortLayers,
        sortFunction: this.sortFunction,
      },
      this,
    );
    this.mapObject = new Control.Layers(null, null, options);

    propsBinder(this, this.mapObject, this.$options.props);

    this.mapObject.setElement(this.$el);
    this.$parent.registerLayerControl(this);

    this.$nextTick(() => {
      /**
       * Triggers when the component is ready
       * @type {object}
       * @property {object} mapObject - reference to leaflet map object
       */
      this.$emit('ready', this.mapObject);

      // eslint-disable-next-line no-underscore-dangle
      this.allLayers = this.mapObject._layers;
    });
  },
  methods: {
    ...LControlLayers.methods,

    isLayerActive(layer) {
      // eslint-disable-next-line no-underscore-dangle
      return this.activeLayersIds.includes(layer._leaflet_id);
    },

    setActiveLayers() {
      this.activeLayersIds = this.allLayers.reduce((acc, { layer }) => {
        if (this.map.hasLayer(layer)) {
          // eslint-disable-next-line no-underscore-dangle
          acc.push(layer._leaflet_id);
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
      this.layers.forEach(({ layer }) => {
        if (this.isLayerActive(layer)) {
          this.removeLayerFromMap(layer);
        }
      });
    },

    enableLayer({ layer }) {
      this.removeActiveLayers();

      this.addLayerToMap(layer);
      this.setActiveLayers();
    },

    enableOverlay(overlay, value) {
      if (value) {
        this.addLayerToMap(overlay.layer);
      } else {
        this.removeLayerFromMap(overlay.layer);
      }

      this.setActiveLayers();
    },
  },
};
</script>

<style lang="scss">
.geomap-layers-control {
  min-width: 200px;
}
</style>
