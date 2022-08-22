<template lang="pug">
  v-expansion-panel(color="grey")
    v-expansion-panel-content
      template(#header="")
        span.white {{ $t('flowchart.properties') }}
      v-divider
      v-card
        v-card-text
          flowchart-color-field.my-1(
            v-if="showFill",
            :label="$t('flowchart.fill')",
            :value="fillValue",
            :palette="shapesColors",
            @input="updateFill"
          )
          flowchart-color-field.my-1(
            v-if="showStroke",
            :label="$t('flowchart.stroke')",
            :value="stroke",
            :palette="borderColors",
            @input="updateStroke"
          )
          template(v-if="showStroke && isStrokeEnabled")
            flowchart-number-field.my-2(
              :label="$t('flowchart.strokeWidth')",
              :value="strokeWidth",
              @input="updateStrokeWidth"
            )
            flowchart-stroke-type-field.my-2(
              :label="$t('flowchart.strokeType')",
              :value="strokeType",
              @input="updateStrokeType"
            )
          flowchart-line-type-field.my-2(
            v-if="showLineType",
            :label="$t('flowchart.lineType')",
            :value="lineType",
            @input="updateLineType"
          )
          v-divider
          flowchart-color-field.my-1(
            :label="$t('flowchart.fontColor')",
            :value="textColor",
            :palette="textColors",
            @input="updateTextColor"
          )
          flowchart-color-field.my-1(
            :label="$t('flowchart.fontBackgroundColor')",
            :value="textBackgroundColor",
            :palette="backgroundColors",
            @input="updateTextBackgroundColor"
          )
          flowchart-number-field.my-2(
            :label="$t('flowchart.fontSize')",
            :value="textFontSize",
            @input="updateFontSize"
          )
</template>

<script>
import { merge, get } from 'lodash';

import { COLORS } from '@/config';

import { STROKE_TYPES } from '@/constants';

import { formBaseMixin } from '@/mixins/form';

import FlowchartColorField from './partials/flowchart-color-field.vue';
import FlowchartNumberField from './partials/flowchart-number-field.vue';
import FlowchartStrokeTypeField from './partials/flowchart-stroke-type-field.vue';
import FlowchartLineTypeField from './partials/flowchart-line-type-field.vue';

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
      return this.getShapesPropertyValue(this.shapesWithStroke, 'properties.fill');
    },

    stroke() {
      return this.getShapesPropertyValue(this.shapesWithStroke, 'properties.stroke');
    },

    isStrokeEnabled() {
      return this.stroke && this.stroke !== 'transparent';
    },

    lineType() {
      return this.getShapesPropertyValue(this.shapesWithStroke, 'lineType');
    },

    strokeWidth() {
      return this.getShapesPropertyValue(this.shapesWithStroke, 'properties.stroke-width');
    },

    strokeType() {
      return this.getShapesPropertyValue(this.shapesWithStroke, 'stroke-dasharray')
        ? STROKE_TYPES.dashed
        : STROKE_TYPES.solid;
    },

    textColor() {
      return this.getShapesPropertyValue(this.selectedShapes, 'textProperties.color');
    },

    textBackgroundColor() {
      return this.getShapesPropertyValue(this.selectedShapes, 'textProperties.backgroundColor');
    },

    textFontSize() {
      return this.getShapesPropertyValue(this.selectedShapes, 'textProperties.fontSize');
    },
  },
  methods: {
    getShapesPropertyValue(shapes, path) {
      const [firstShape] = shapes;
      const value = get(firstShape, path);

      return shapes?.every(shape => get(shape, path) === value)
        ? value
        : undefined;
    },

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
        'stroke-dasharray': strokeType === STROKE_TYPES.solid ? undefined : '4 4',
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
