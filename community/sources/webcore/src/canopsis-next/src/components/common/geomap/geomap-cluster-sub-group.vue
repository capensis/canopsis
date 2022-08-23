<template lang="pug">
  div(style="display: none;")
    slot(v-if="ready")
</template>

<script>
import { LayerGroupMixin, optionsMerger, findRealParent, propsBinder } from 'vue2-leaflet';
import { FeatureGroup, DomEvent } from 'leaflet';

import 'leaflet.featuregroup.subgroup';

export default {
  mixins: [LayerGroupMixin],
  data() {
    return {
      ready: false,
    };
  },
  mounted() {
    this.parentContainer = findRealParent(this.$parent, true);

    const options = optionsMerger(
      {},
      this,
    );
    this.mapObject = new FeatureGroup.SubGroup(this.parentContainer, options);

    propsBinder(this, this.mapObject, this.$options.props);
    DomEvent.on(this.mapObject, this.$listeners);

    this.ready = true;

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
