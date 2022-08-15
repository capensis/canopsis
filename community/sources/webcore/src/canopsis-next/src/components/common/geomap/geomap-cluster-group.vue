<template lang="pug">
  div(style="display: none;")
    slot(v-if="ready")
</template>

<script>
import { MarkerClusterGroup } from 'leaflet.markercluster';
import { LayerGroupMixin, findRealParent, propsBinder } from 'vue2-leaflet';
import { DomEvent } from 'leaflet';

export default {
  mixins: [LayerGroupMixin],
  props: {
    polygonOptions: {
      type: Object,
      default: () => ({
        color: '#5a6D80',
      }),
    },
  },
  data() {
    return {
      ready: false,
    };
  },
  mounted() {
    this.mapObject = new MarkerClusterGroup(this.$options.props);
    propsBinder(this, this.mapObject, this.$options.props);
    DomEvent.on(this.mapObject, this.$listeners);
    this.ready = true;
    this.parentContainer = findRealParent(this.$parent, true);
    if (this.visible) {
      this.parentContainer.addLayer(this);
    }
    this.$nextTick(() => {
      /**
       * Triggers when the component is ready
       * @type {object}
       * @property {object} mapObject - reference to leaflet map object
       */
      this.$emit('ready', this.mapObject);
    });
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

  background: #5a6D80;
  border: 1px solid #FFFFFF;
  border-radius: 50%;
}
</style>
