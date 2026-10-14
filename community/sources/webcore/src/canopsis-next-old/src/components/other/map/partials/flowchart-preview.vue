<template>
  <div class="flowchart-preview">
    <flowchart
      :shapes="shapes"
      :background-color="map.parameters.background_color"
      readonly
    >
      <template #layers="{ data }">
        <flowchart-points-preview
          :points="map.parameters.points"
          :popup-template="popupTemplate"
          :popup-actions="popupActions"
          :color-indicator="colorIndicator"
          :pbehavior-enabled="pbehaviorEnabled"
          :shapes="data"
          v-on="$listeners"
        />
      </template>
    </flowchart>
    <c-help-icon
      :text="$t('flowchart.panzoom.helpText')"
      size="32"
      icon-class="map-preview__help-icon"
      color="secondary"
      icon="help"
      top
    />
  </div>
</template>

<script>
import { keyBy } from 'lodash';

import { getCSSVariableColor, getCSSVariableName, getDarkenColor, isCSSVariable } from '@/helpers/color';
import { getEntityColor } from '@/helpers/entities/entity/color';
import { isNotActivePbehaviorType } from '@/helpers/entities/pbehavior/form';

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
      const entityColor = getEntityColor(point.entity, this.colorIndicator);
      const color = isCSSVariable(entityColor)
        ? getCSSVariableColor(document.body, getCSSVariableName(entityColor))
        : entityColor;
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
  border-radius: 5px;
  overflow: hidden;
}
</style>
