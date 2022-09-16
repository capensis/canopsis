<template lang="pug">
  div.flowchart-preview
    flowchart(
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
    c-help-icon.flowchart-preview__help-icon(size="32", color="secondary", icon="help", top)
      div.pre-wrap(v-html="$t('flowchart.panzoom.helpText')")
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
  position: relative;
  height: 800px;
  width: 100%;

  &__help-icon {
    cursor: pointer;
    position: absolute;
    bottom: 10px;
    right: 10px;
  }
}
</style>
