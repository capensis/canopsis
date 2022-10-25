<template lang="pug">
  div.flowchart-preview
    flowchart(
      :shapes="shapes",
      :background-color="map.parameters.background_color",
      readonly
    )
      template(#layers="{ data }")
        flowchart-points-preview(
          v-on="$listeners",
          :points="map.parameters.points",
          :popup-template="popupTemplate",
          :popup-actions="popupActions",
          :color-indicator="colorIndicator",
          :pbehavior-enabled="pbehaviorEnabled",
          :shapes="data"
        )
    c-help-icon.map-preview__help-icon(size="32", color="secondary", icon="help", top)
      div.pre-wrap(v-html="$t('flowchart.panzoom.helpText')")
</template>

<script>
import { keyBy } from 'lodash';

import { getDarkenColor, getEntityColor } from '@/helpers/color';
import { isNotActivePbehaviorType } from '@/helpers/entities/pbehavior';

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
    pointsByShape() {
      return keyBy(this.map.parameters.points, 'shape');
    },

    shapes() {
      if (this.colorIndicator || this.pbehaviorEnabled) {
        return this.map.parameters.shapes.reduce((acc, shape) => {
          const point = this.pointsByShape[shape._id];

          acc[shape._id] = point
            ? this.getShapeByPoint(shape, point)
            : shape;

          return acc;
        }, {});
      }

      return keyBy(this.map.parameters.shapes, '_id');
    },
  },
  methods: {
    getShapeByEntity(shape, point) {
      const color = getEntityColor(point.entity, this.colorIndicator);
      const darkenColor = getDarkenColor(color, 20);

      return {
        ...shape,
        properties: {
          ...shape.properties,
          fill: color,
          stroke: darkenColor,
        },
        textProperties: {
          ...shape.textProperties,
          color: darkenColor,
        },
      };
    },

    getShapeByPoint(shape, point) {
      if (this.pbehaviorEnabled && isNotActivePbehaviorType(point.entity?.pbehavior_info?.canonical_type)) {
        return shape;
      }

      return this.getShapeByEntity(shape, point);
    },
  },
};
</script>

<style lang="scss">
.flowchart-preview {
  position: relative;
  height: 800px;
  width: 100%;
}
</style>