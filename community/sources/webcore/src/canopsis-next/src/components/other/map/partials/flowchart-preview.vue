<template lang="pug">
  flowchart.flowchart-preview(
    :shapes="shapes",
    :background-color="map.parameters.background_color",
    readonly
  )
    template(#layers="{ data }")
      flowchart-points-preview(
        :points="map.parameters.points",
        :popup-template="popupTemplate",
        :popup-actions="popupActions",
        :color-indicator="colorIndicator",
        :pbehavior-enabled="pbehaviorEnabled",
        :shapes="data"
      )
</template>

<script>
import { keyBy } from 'lodash';

import Flowchart from '@/components/common/flowchart/flowchart.vue';

import FlowchartPointsPreview from './flowchart-points-preview.vue';

export default {
  components: { Flowchart, FlowchartPointsPreview },
  props: {
    map: {
      type: Object,
      required: true,
    },
    popupTemplate: {
      type: String,
      required: false,
    },
    popupActions: {
      type: Boolean,
      default: false,
    },
    colorIndicator: {
      type: String,
      required: false,
    },
    pbehaviorEnabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    shapes() {
      return keyBy(this.map.parameters.shapes, '_id');
    },
  },
};
</script>

<style lang="scss">
.flowchart-preview {
  height: 800px;
  width: 100%;
}
</style>
