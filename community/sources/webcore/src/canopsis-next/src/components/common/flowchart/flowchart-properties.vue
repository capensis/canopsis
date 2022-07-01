<template lang="pug">
  v-expansion-panel(color="grey")
    v-expansion-panel-content
      template(#header="")
        span.white {{ $t('flowchart.properties') }}
      v-card
        v-card-text
          flowchart-color-field.my-1(
            v-if="shapesWithFill.length",
            :label="$t('flowchart.fill')",
            :value="properties.fill",
            @input="updateFill"
          )
          flowchart-color-field.my-1(
            v-if="shapesWithStroke.length",
            :label="$t('flowchart.stroke')",
            :value="properties.stroke",
            @input="updateStroke"
          )
          flowchart-number-field.my-2(
            v-if="shapesWithStroke.length",
            :label="$t('flowchart.strokeWidth')",
            :value="properties.strokeWidth",
            @input="updateStrokeWidth"
          )
          v-divider
          flowchart-color-field.my-1(
            :label="$t('flowchart.fontColor')",
            :value="textProperties.color",
            @input="updateTextColor"
          )
          flowchart-color-field.my-1(
            :label="$t('flowchart.fontBackgroundColor')",
            :value="textProperties.backgroundColor",
            @input="updateTextBackgroundColor"
          )
          flowchart-number-field.my-2(
            :label="$t('flowchart.fontSize')",
            :value="textProperties.fontSize",
            @input="updateFontSize"
          )
</template>

<script>
import { formBaseMixin } from '@/mixins/form';

import FlowchartColorField from './partials/flowchart-color-field.vue';
import FlowchartNumberField from './partials/flowchart-number-field.vue';

export default {
  components: {
    FlowchartColorField,
    FlowchartNumberField,
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
    selectedShapes() {
      return this.selected.map(id => this.shapes[id]);
    },

    shapesWithFill() {
      return this.selectedShapes.filter(({ properties }) => properties.fill);
    },

    shapesWithStroke() {
      return this.selectedShapes.filter(({ properties }) => properties.stroke);
    },

    properties() {
      const [firstShapeWithFill] = this.shapesWithFill;
      const [firstShapeWithStroke] = this.shapesWithStroke;

      return {
        fill: firstShapeWithFill?.properties?.fill,
        stroke: firstShapeWithStroke?.properties?.stroke,
        strokeWidth: firstShapeWithStroke?.properties?.['stroke-width'],
      };
    },

    textProperties() {
      const [firstShape] = this.selectedShapes;

      return firstShape.textProperties;
    },
  },
  methods: {
    updateSelectedShapesProperties(properties) {
      this.updateModel(this.selected.reduce((acc, id) => {
        const shape = this.shapes[id];

        acc[id] = {
          ...shape,
          properties: {
            ...shape.properties,
            ...properties,
          },
        };

        return acc;
      }, { ...this.shapes }));
    },

    updateSelectedShapesTextProperties(textProperties) {
      this.updateModel(this.selected.reduce((acc, id) => {
        const shape = this.shapes[id];

        acc[id] = {
          ...shape,
          textProperties: {
            ...shape.textProperties,
            ...textProperties,
          },
        };

        return acc;
      }, { ...this.shapes }));
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
