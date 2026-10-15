<template>
  <div style="display: none">
    <slot v-if="ready" />
  </div>
</template>

<script>
import * as L from 'leaflet';

import { GEOMAP_PANES } from '@/constants';

import 'leaflet.markercluster/dist/MarkerCluster.css';
import 'leaflet.markercluster/dist/MarkerCluster.Default.css';

const CustomDivIcon = L.DivIcon.extend({
  createIcon(oldIcon) {
    const div = L.DivIcon.prototype.createIcon.apply(this, oldIcon);

    if (this.options.style) {
      Object.entries(this.options.style).forEach(([name, value]) => {
        div.style[name] = value;
      });
    }

    return div;
  },
});

export default {
  inject: ['geomapRoot'],

  provide() {
    return {
      geomapRoot: this.geomapRoot,
      geomapParent: this,
    };
  },

  props: {
    name: {
      type: String,
      required: false,
    },
    layerType: {
      type: String,
      default: 'overlay',
    },
    showCoverageOnHover: {
      type: Boolean,
      default: false,
    },
    clusterClassName: {
      type: String,
      required: false,
    },
    clusterStyle: {
      type: Object,
      required: false,
    },
  },

  emits: ['ready'],

  data() {
    return {
      mapObject: undefined,
      ready: false,
    };
  },

  watch: {
    clusterClassName() {
      this.mapObject.refreshClusters();
    },
    clusterStyle() {
      this.mapObject.refreshClusters();
    },
  },

  async mounted() {
    window.L = L;
    await import('leaflet.markercluster');

    this.mapObject = new L.MarkerClusterGroup({
      iconCreateFunction: this.createIcon,
      showCoverageOnHover: this.showCoverageOnHover,
      clusterPane: GEOMAP_PANES.clusters,
    });

    this.ready = true;
    this.geomapRoot.addLayer(this);

    this.$nextTick(() => {
      this.$emit('ready', this.mapObject);
    });
  },

  beforeUnmount() {
    this.geomapRoot.removeLayer(this);
  },

  methods: {
    addLayer(layer) {
      this.mapObject.addLayer(layer.mapObject);
    },

    removeLayer(layer) {
      this.mapObject.removeLayer(layer.mapObject);
    },

    createIcon(cluster) {
      const count = cluster.getChildCount();

      return new CustomDivIcon({
        html: `<div><span>${count}</span></div>`,
        className: `marker-cluster ${this.clusterClassName || ''}`,
        style: this.clusterStyle,
        iconSize: new L.Point(40, 40),
      });
    },
  },
};
</script>

<style lang="scss">
.marker-cluster {
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  width: 30px;
  height: 30px;

  background: #5a6d80;
  border: 1px solid #ffffff;
  border-radius: 50%;

  div {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    height: 100%;
    margin: 0;
    border-radius: 50%;
  }

  span {
    line-height: 1;
  }
}
</style>
