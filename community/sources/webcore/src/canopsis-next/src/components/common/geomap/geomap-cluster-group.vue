<template>
  <div style="display: none;">
    <slot v-if="ready" />
  </div>
</template>

<script>
import { MarkerClusterGroup } from 'leaflet.markercluster';
import { LayerGroupMixin, optionsMerger, findRealParent, propsBinder } from 'vue2-leaflet';
import { DivIcon, Point, DomEvent } from 'leaflet';

const CustomDivIcon = DivIcon.extend({
  createIcon(oldIcon) {
    const div = DivIcon.prototype.createIcon.apply(this, oldIcon);

    if (this.options.style) {
      Object.entries(this.options.style).forEach(([name, value]) => {
        div.style[name] = value;
      });
    }

    return div;
  },
});

export default {
  mixins: [LayerGroupMixin],
  props: {
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
  data() {
    return {
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
  mounted() {
    const options = optionsMerger(
      {
        iconCreateFunction: this.createIcon,
        showCoverageOnHover: this.showCoverageOnHover,
      },
      this,
    );
    this.mapObject = new MarkerClusterGroup(options);

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
  methods: {
    createIcon(cluster) {
      const count = cluster.getChildCount();

      return new CustomDivIcon({
        html: `<div><span>${count}</span></div>`,
        className: `marker-cluster ${this.clusterClassName}`,
        style: this.clusterStyle,
        iconSize: new Point(40, 40),
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

  background: #5a6D80;
  border: 1px solid #FFFFFF;
  border-radius: 50%;
}
</style>
