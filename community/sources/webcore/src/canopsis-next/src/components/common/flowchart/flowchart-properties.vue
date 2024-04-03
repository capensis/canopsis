<template>
  <v-expansion-panels color="grey">
    <v-expansion-panel color="grey">
      <v-expansion-panel-header>
        <span>{{ $t('flowchart.properties') }}</span>
      </v-expansion-panel-header>
      <v-expansion-panel-content>
        <v-divider />
        <v-card flat>
          <v-card-text>
            <flowchart-color-field
              v-if="showFill"
              :label="$t('flowchart.fill')"
              :value="fillValue"
              :palette="shapesColors"
              class="my-1"
              @input="updateFill"
            />
            <flowchart-color-field
              v-if="showStroke"
              :label="$t('flowchart.stroke')"
              :value="stroke"
              :palette="borderColors"
              class="my-1"
              @input="updateStroke"
            />
            <template v-if="showStroke && isStrokeEnabled">
              <flowchart-number-field
                :label="$t('flowchart.strokeWidth')"
                :value="strokeWidth"
                class="my-1"
                @input="updateStrokeWidth"
              />
              <flowchart-stroke-type-field
                :label="$t('flowchart.strokeType')"
                :value="strokeType"
                class="my-1"
                @input="updateStrokeType"
              />
            </template>
            <flowchart-line-type-field
              v-if="showLineType"
              :label="$t('flowchart.lineType')"
              :value="lineType"
              :average-points="lineAveragePoints"
              class="my-2"
              @input="updateLineType"
            />
            <v-divider v-if="showLineType || showStroke || showFill" />
            <flowchart-color-field
              :label="$t('flowchart.fontColor')"
              :value="textColor"
              :palette="textColors"
              class="my-1"
              @input="updateTextColor"
            />
            <flowchart-color-field
              :label="$t('flowchart.fontBackgroundColor')"
              :value="textBackgroundColor"
              :palette="backgroundColors"
              class="my-1"
              @input="updateTextBackgroundColor"
            />
            <flowchart-number-field
              :label="$t('flowchart.fontSize')"
              :value="textFontSize"
              class="my-1"
              @input="updateFontSize"
            />
          </v-card-text>
        </v-card>
      </v-expansion-panel-content>
    </v-expansion-panel>
  </v-expansion-panels>
</template>

<script>
import { merge } from 'lodash';

import { COLORS } from '@/config';
import { STROKE_TYPES } from '@/constants';

import { getPropertyValueByShapes } from '@/helpers/flowchart/shapes';
import { calculateCenterBetweenPoint } from '@/helpers/flowchart/points';

import { formBaseMixin } from '@/mixins/form';

import FlowchartColorField from './fields/flowchart-color-field.vue';
import FlowchartNumberField from './fields/flowchart-number-field.vue';
import FlowchartStrokeTypeField from './fields/flowchart-stroke-type-field.vue';
import FlowchartLineTypeField from './fields/flowchart-line-type-field.vue';

export default {
  components: {
    FlowchartColorField,
    FlowchartNumberField,
    FlowchartStrokeTypeField,
    FlowchartLineTypeField,
  },
  mixins: [formBaseMixin],
  model: {
    prop: 'shapes',
    event: 'input',
  },
  props: {
    shapes: {
      type: Object,
      required: true,
    },
    selected: {
      type: Array,
      required: true,
    },
  },
  computed: {
    shapesColors() {
      return COLORS.flowchart.shapes;
    },

    borderColors() {
      return COLORS.flowchart.border;
    },

    textColors() {
      return COLORS.flowchart.text;
    },

    backgroundColors() {
      return COLORS.flowchart.background;
    },

    selectedShapes() {
      return this.selected.map(id => this.shapes[id]);
    },

    shapesWithFill() {
      return this.selectedShapes.filter(({ properties }) => properties.fill);
    },

    shapesWithStroke() {
      return this.selectedShapes.filter(({ properties }) => properties.stroke);
    },

    shapesWithLineType() {
      return this.selectedShapes.filter(({ lineType }) => lineType);
    },

    showFill() {
      return this.selected.length === this.shapesWithFill.length;
    },

    showStroke() {
      return this.selected.length === this.shapesWithStroke.length;
    },

    showLineType() {
      return this.selected.length === this.shapesWithLineType.length;
    },

    fillValue() {
      return getPropertyValueByShapes(this.shapesWithFill, 'properties.fill');
    },

    stroke() {
      return getPropertyValueByShapes(this.shapesWithStroke, 'properties.stroke');
    },

    isStrokeEnabled() {
      return this.stroke && this.stroke !== 'transparent';
    },

    lineType() {
      return getPropertyValueByShapes(this.shapesWithLineType, 'lineType');
    },

    strokeWidth() {
      return getPropertyValueByShapes(this.shapesWithStroke, 'properties.stroke-width');
    },

    strokeType() {
      return getPropertyValueByShapes(this.shapesWithStroke, 'properties.stroke-dasharray')
        ? STROKE_TYPES.dashed
        : STROKE_TYPES.solid;
    },

    textColor() {
      return getPropertyValueByShapes(this.selectedShapes, 'textProperties.color');
    },

    textBackgroundColor() {
      return getPropertyValueByShapes(this.selectedShapes, 'textProperties.backgroundColor');
    },

    textFontSize() {
      return getPropertyValueByShapes(this.selectedShapes, 'textProperties.fontSize');
    },

    lineAveragePoints() {
      return this.shapesWithLineType.reduce((acc, { points }) => {
        if (acc.length) {
          return acc.map((point, index) => calculateCenterBetweenPoint(
            acc[index],
            point,
          ));
        }

        return [...points];
      }, []);
    },
  },
  methods: {
    updateSelectedShapes(shapeData) {
      this.updateModel(this.selected.reduce((acc, id) => {
        const shape = this.shapes[id];

        acc[id] = merge({}, shape, shapeData);

        return acc;
      }, { ...this.shapes }));
    },

    updateSelectedShapesProperties(properties) {
      this.updateSelectedShapes({ properties });
    },

    updateSelectedShapesTextProperties(textProperties) {
      this.updateSelectedShapes({ textProperties });
    },

    updateFill(fill) {
      this.updateSelectedShapesProperties({ fill });
    },

    updateStroke(stroke) {
      this.updateSelectedShapesProperties({ stroke });
    },

    updateStrokeWidth(strokeWidth) {
      this.updateSelectedShapesProperties({ 'stroke-width': strokeWidth });
    },

    updateStrokeType(strokeType) {
      this.updateSelectedShapesProperties({
        'stroke-dasharray': strokeType === STROKE_TYPES.solid ? '' : '4 4',
      });
    },

    updateLineType(lineType) {
      this.updateSelectedShapes({ lineType });
    },

    updateTextColor(color) {
      this.updateSelectedShapesTextProperties({ color });
    },

    updateTextBackgroundColor(backgroundColor) {
      this.updateSelectedShapesTextProperties({ backgroundColor });
    },

    updateFontSize(fontSize) {
      this.updateSelectedShapesTextProperties({ fontSize });
    },
  },
};
</script>
