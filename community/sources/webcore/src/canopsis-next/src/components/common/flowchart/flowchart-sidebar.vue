<template lang="pug">
  v-navigation-drawer(permanent, width="300")
    v-expansion-panel(color="grey")
      v-expansion-panel-content
        template(#header="")
          span.white {{ $t('flowchart.shapes') }}
        v-layout(row, wrap)
          v-btn.ma-0.pa-0(v-for="button in buttons", :key="button.icon", flat, large, @click="button.action")
            v-icon(size="38", color="grey darken-3") {{ button.icon }}
          file-selector(ref="fileSelector", hide-details, @change="addImage")
            template(#activator="{ on }")
              v-btn.ma-0.pa-0(v-on="on", flat, large)
                v-icon(size="38", color="grey darken-3") $vuetify.icons.image_shape
</template>

<script>
import { getFileDataUrlContent } from '@/helpers/file/file-select';
import {
  generateArrowLineShape,
  generateBidirectionalArrowLineShape,
  generateCircleShape,
  generateEllipseShape,
  generateImageShape,
  generateLineShape,
  generateParallelogramShape,
  generateRectShape,
  generateRhombusShape,
  generateStorageShape,
} from '@/helpers/flowchart/shapes';

import { formBaseMixin } from '@/mixins/form';

import FileSelector from '@/components/forms/fields/file-selector.vue';
import { generatePoint } from '@/helpers/flowchart/points';
import { getImageProperties } from '@/helpers/file/image';

export default {
  components: { FileSelector },
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
    viewBox: {
      type: Object,
      required: true,
    },
  },
  computed: {
    buttons() {
      return [
        { icon: '$vuetify.icons.rect_shape', action: this.addRectangle },
        { icon: '$vuetify.icons.rounded_rect_shape', action: this.addRoundedRectangle },
        { icon: '$vuetify.icons.square_shape', action: this.addSquare },
        { icon: '$vuetify.icons.rhombus_shape', action: this.addRhombus },
        { icon: '$vuetify.icons.circle_shape', action: this.addCircle },
        { icon: '$vuetify.icons.ellipse_shape', action: this.addEllipse },
        { icon: '$vuetify.icons.parallelogram_shape', action: this.addParallelogram },
        { icon: '$vuetify.icons.storage_shape', action: this.addStorage },
        { icon: '$vuetify.icons.line_shape', action: this.addLine },
        { icon: '$vuetify.icons.dashed_line_shape', action: this.addDashedLine },
        { icon: '$vuetify.icons.arrow_line_shape', action: this.addArrowLine },
        { icon: '$vuetify.icons.dashed_arrow_line_shape', action: this.addDashedArrowLine },
        { icon: '$vuetify.icons.bidirectional_arrow_line_shape', action: this.addBidirectionalArrowLine },
        { icon: '$vuetify.icons.dashed_bidirectional_arrow_line_shape', action: this.addDashedBidirectionalArrowLine },
        { icon: '$vuetify.icons.text_shape', action: this.addText },
        { icon: '$vuetify.icons.textbox_shape', action: this.addTextbox },
      ];
    },

    viewBoxCenter() {
      return {
        x: this.viewBox.x + this.viewBox.width / 2,
        y: this.viewBox.y + this.viewBox.height / 2,
      };
    },

    centerRectProperties() {
      const width = 100;
      const height = 100;

      return {
        x: this.viewBoxCenter.x - width / 2,
        y: this.viewBoxCenter.y - height / 2,
        width,
        height,
      };
    },

    centerCircleProperties() {
      const diameter = 100;
      const halfDiameter = diameter / 2;

      return {
        x: this.viewBoxCenter.x - halfDiameter,
        y: this.viewBoxCenter.y - halfDiameter,
        diameter,
      };
    },

    centerLinePoints() {
      const length = 100;
      const halfLength = length / 2;

      return [
        generatePoint({
          x: this.viewBoxCenter.x - halfLength,
          y: this.viewBoxCenter.y,
        }),
        generatePoint({
          x: this.viewBoxCenter.x + halfLength,
          y: this.viewBoxCenter.y,
        }),
      ];
    },
  },
  methods: {
    addShape(shape) {
      this.updateModel({
        ...this.shapes,
        [shape._id]: shape,
      });
    },

    addRectangle() {
      const rect = generateRectShape({
        ...this.centerRectProperties,
        text: 'Rectangle',
        textProperties: {
          alignCenter: true,
          justifyCenter: true,
        },
        properties: {
          stroke: 'black',
          'stroke-width': 2,
          fill: 'white',
        },
      });

      this.addShape(rect);
    },

    addRoundedRectangle() {
      const rect = generateRectShape({
        ...this.centerRectProperties,
        text: 'Rounded rectangle',
        textProperties: {
          alignCenter: true,
          justifyCenter: true,
        },
        properties: {
          rx: 20,
          ry: 20,
          stroke: 'black',
          'stroke-width': 2,
          fill: 'white',
        },
      });

      this.addShape(rect);
    },

    addLine() {
      const line = generateLineShape({
        points: this.centerLinePoints,
        text: 'Line',
        properties: {
          stroke: 'black',
          'stroke-width': 2,
        },
      });

      this.addShape(line);
    },

    addDashedLine() {
      const line = generateLineShape({
        points: this.centerLinePoints,
        text: 'Dashed line',
        properties: {
          stroke: 'black',
          'stroke-width': 2,
        },
      });

      this.addShape(line);
    },

    addDashedArrowLine() {
      const arrowLine = generateArrowLineShape({
        points: this.centerLinePoints,
        text: 'Dashed arrow line',
        properties: {
          stroke: 'black',
          'stroke-width': 2,
          'stroke-dasharray': '4 4',
        },
      });

      this.addShape(arrowLine);
    },

    addArrowLine() {
      const arrowLine = generateArrowLineShape({
        points: this.centerLinePoints,
        text: 'Arrow line',
        properties: {
          stroke: 'black',
          'stroke-width': 2,
        },
      });

      this.addShape(arrowLine);
    },

    addBidirectionalArrowLine() {
      const bidirectionalArrowLine = generateBidirectionalArrowLineShape({
        points: this.centerLinePoints,
        text: 'Bidirectional arrow line',
        properties: {
          stroke: 'black',
          'stroke-width': 2,
        },
      });

      this.addShape(bidirectionalArrowLine);
    },

    addDashedBidirectionalArrowLine() {
      const bidirectionalArrowLine = generateBidirectionalArrowLineShape({
        points: this.centerLinePoints,
        text: 'Dashed bidirectional arrow line',
        properties: {
          stroke: 'black',
          'stroke-width': 2,
          'stroke-dasharray': '4 4',
        },
      });

      this.addShape(bidirectionalArrowLine);
    },

    addCircle() {
      const circle = generateCircleShape({
        ...this.centerCircleProperties,
        text: 'Circle',
        textProperties: {
          alignCenter: true,
          justifyCenter: true,
        },
        properties: {
          stroke: 'black',
          'stroke-width': 2,
          fill: 'white',
        },
      });

      this.addShape(circle);
    },

    addEllipse() {
      const ellipse = generateEllipseShape({
        ...this.centerRectProperties,
        text: 'Ellipse',
        textProperties: {
          alignCenter: true,
          justifyCenter: true,
        },
        properties: {
          stroke: 'black',
          'stroke-width': 2,
          fill: 'white',
        },
      });

      this.addShape(ellipse);
    },

    addRhombus() {
      const rhombus = generateRhombusShape({
        ...this.centerRectProperties,
        text: 'Rhombus',
        textProperties: {
          alignCenter: true,
          justifyCenter: true,
        },
        properties: {
          stroke: 'black',
          'stroke-width': 2,
          fill: 'white',
        },
      });

      this.addShape(rhombus);
    },

    addParallelogram() {
      const parallelogram = generateParallelogramShape({
        ...this.centerRectProperties,
        text: 'Parallelogram',
        textProperties: {
          alignCenter: true,
          justifyCenter: true,
        },
        properties: {
          stroke: 'black',
          'stroke-width': 2,
          fill: 'white',
        },
      });

      this.addShape(parallelogram);
    },

    addSquare() {
      const square = generateRectShape({
        ...this.centerRectProperties,
        text: 'Square',
        textProperties: {
          alignCenter: true,
          justifyCenter: true,
        },
        properties: {
          stroke: 'black',
          'stroke-width': 2,
          fill: 'white',
        },
        aspectRatio: true,
      });

      this.addShape(square);
    },

    addText() {
      const text = generateRectShape({
        ...this.centerRectProperties,
        text: 'Text',
        textProperties: {
          alignCenter: true,
          justifyCenter: true,
        },
        properties: {
          fill: 'transparent',
        },
      });

      this.addShape(text);
    },

    addTextbox() {
      const textbox = generateRectShape({
        ...this.centerRectProperties,
        text: '<h2>Heading</h2><p>Paragraph</p>',
        textProperties: {
          hidden: true,
        },
        properties: {
          fill: 'transparent',
        },
      });

      this.addShape(textbox);
    },

    addStorage() {
      const storage = generateStorageShape({
        ...this.centerRectProperties,
        text: 'Storage',
        textProperties: {
          alignCenter: true,
          justifyCenter: true,
        },
        properties: {
          stroke: 'black',
          'stroke-width': 2,
          fill: 'grey',
        },
      });

      this.addShape(storage);
    },

    async addImage([file]) {
      const src = await getFileDataUrlContent(file);
      const { width, height } = await getImageProperties(src);

      const image = generateImageShape({
        ...this.centerRectProperties,
        x: this.viewBoxCenter.x - width / 2,
        y: this.viewBoxCenter.y - height / 2,
        width,
        height,
        src,
        text: file.name,
        aspectRatio: true,
      });

      this.addShape(image);

      this.$refs.fileSelector.clear();
    },
  },
};
</script>
