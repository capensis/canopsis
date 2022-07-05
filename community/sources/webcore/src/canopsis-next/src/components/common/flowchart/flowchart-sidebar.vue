<template lang="pug">
  v-navigation-drawer.flowchart-sidebar(permanent, width="300")
    v-expansion-panel(color="grey", expand)
      v-expansion-panel-content
        template(#header="")
          span.white {{ $t('flowchart.shapes') }}
        v-divider
        v-layout(row, wrap)
          v-btn.ma-0.pa-0.flowchart-sidebar__button(
            v-for="(button, index) in buttons",
            :key="index",
            flat,
            large,
            @click="button.action"
          )
            component.grey--text.text--darken-3.pa-1(:is="button.icon")
          file-selector(ref="fileSelector", hide-details, @change="addImage")
            template(#activator="{ on }")
              v-btn.ma-0.pa-0.flowchart-sidebar__button(v-on="on", flat, large)
                image-shape-icon.grey--text.text--darken-3
      v-expansion-panel-content
        template(#header="")
          span.white {{ $t('flowchart.icons') }}
        v-divider
        v-layout(row, wrap)
          v-btn.ma-0(v-for="icon in icons", :key="icon.src", flat, fab, small, @click="addIconAsset(icon.src)")
            img(:src="icon.src")
</template>

<script>
import { LINE_TYPES } from '@/constants';

import { getFileDataUrlContent } from '@/helpers/file/file-select';
import {
  generateArrowLineShape,
  generateBidirectionalArrowLineShape,
  generateCircleShape,
  generateEllipseShape,
  generateImageShape,
  generateLineShape,
  generateParallelogramShape,
  generateProcessShape,
  generateDocumentShape,
  generateRectShape,
  generateRhombusShape,
  generateStorageShape,
} from '@/helpers/flowchart/shapes';
import { generatePoint } from '@/helpers/flowchart/points';
import { getImageProperties } from '@/helpers/file/image';

import { formBaseMixin } from '@/mixins/form';

import FileSelector from '@/components/forms/fields/file-selector.vue';

import RectShapeIcon from './icons/rect-shape.vue';
import RoundedRectShapeIcon from './icons/rounded-rect-shape.vue';
import SquareShapeIcon from './icons/square-shape.vue';
import RhombusShapeIcon from './icons/rhombus-shape.vue';
import CircleShapeIcon from './icons/circle-shape.vue';
import EllipseShapeIcon from './icons/ellipse-shape.vue';
import ParallelogramShapeIcon from './icons/parallelogram-shape.vue';
import StorageShapeIcon from './icons/storage-shape.vue';
import LineShapeIcon from './icons/line-shape.vue';
import ArrowLineShapeIcon from './icons/arrow-line-shape.vue';
import BidirectionalArrowLineShapeIcon from './icons/bidirectional-arrow-line-shape.vue';
import ImageShapeIcon from './icons/image-shape.vue';
import CurveLineShapeIcon from './icons/curve-line-shape.vue';
import CurveArrowLineShapeIcon from './icons/curve-arrow-line-shape.vue';
import BidirectionalCurveArrowLineShape from './icons/bidirectional-curve-arrow-line-shape.vue';
import ProcessShapeIcon from './icons/process-shape.vue';
import DocumentShapeIcon from './icons/document-shape.vue';
import TextShapeIcon from './icons/text-shape.vue';
import TextboxShapeIcon from './icons/textbox-shape.vue';

import assets from './assets';

export default {
  components: { FileSelector, ImageShapeIcon },
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
    selected: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    buttons() {
      return [
        { icon: RectShapeIcon, action: this.addRectangle },
        { icon: RoundedRectShapeIcon, action: this.addRoundedRectangle },
        { icon: SquareShapeIcon, action: this.addSquare },
        { icon: RhombusShapeIcon, action: this.addRhombus },
        { icon: CircleShapeIcon, action: this.addCircle },
        { icon: EllipseShapeIcon, action: this.addEllipse },
        { icon: ParallelogramShapeIcon, action: this.addParallelogram },
        { icon: ProcessShapeIcon, action: this.addProcess },
        { icon: DocumentShapeIcon, action: this.addDocument },
        { icon: StorageShapeIcon, action: this.addStorage },
        { icon: CurveLineShapeIcon, action: this.addCurveLine },
        { icon: CurveArrowLineShapeIcon, action: this.addCurveArrowLine },
        { icon: BidirectionalCurveArrowLineShape, action: this.addBidirectionalCurveArrowLine },
        { icon: LineShapeIcon, action: this.addLine },
        { icon: ArrowLineShapeIcon, action: this.addArrowLine },
        { icon: BidirectionalArrowLineShapeIcon, action: this.addBidirectionalArrowLine },
        { icon: TextShapeIcon, action: this.addText },
        { icon: TextboxShapeIcon, action: this.addTextbox },
      ];
    },

    icons() {
      return assets.map(assetPath => ({
        src: assetPath,
      }));
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

    centerCurveLinePoints() {
      const size = 100;
      const halfSize = size / 2;

      return [
        generatePoint({
          x: this.viewBoxCenter.x - halfSize,
          y: this.viewBoxCenter.y + halfSize,
        }),
        generatePoint({
          x: this.viewBoxCenter.x + halfSize,
          y: this.viewBoxCenter.y + halfSize,
        }),
        generatePoint({
          x: this.viewBoxCenter.x - halfSize,
          y: this.viewBoxCenter.y - halfSize,
        }),
        generatePoint({
          x: this.viewBoxCenter.x + halfSize,
          y: this.viewBoxCenter.y - halfSize,
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

      this.$emit('update:selected', [shape._id]);
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

    addCurveLine() {
      const line = generateLineShape({
        points: this.centerCurveLinePoints,
        lineType: LINE_TYPES.curve,
        text: 'Curve line',
        properties: {
          stroke: 'black',
          'stroke-width': 2,
        },
      });

      this.addShape(line);
    },

    addCurveArrowLine() {
      const arrowLine = generateArrowLineShape({
        points: this.centerCurveLinePoints,
        lineType: LINE_TYPES.curve,
        text: 'Curve arrow line',
        properties: {
          stroke: 'black',
          'stroke-width': 2,
        },
      });

      this.addShape(arrowLine);
    },

    addBidirectionalCurveArrowLine() {
      const bidirectionalArrowLine = generateBidirectionalArrowLineShape({
        points: this.centerCurveLinePoints,
        lineType: LINE_TYPES.curve,
        text: 'Bidirectional arrow line',
        properties: {
          stroke: 'black',
          'stroke-width': 2,
        },
      });

      this.addShape(bidirectionalArrowLine);
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

    addProcess() {
      const parallelogram = generateProcessShape({
        ...this.centerRectProperties,
        text: 'Process',
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

    addDocument() {
      const parallelogram = generateDocumentShape({
        ...this.centerRectProperties,
        text: 'Document',
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
        properties: {
          fill: 'transparent',
          stroke: 'transparent',
        },
      });

      this.addShape(image);

      this.$refs.fileSelector.clear();
    },

    async addIconAsset(src) {
      const { width, height } = await getImageProperties(src);

      const image = generateImageShape({
        ...this.centerRectProperties,
        width,
        height,
        src,
        aspectRatio: true,
        properties: {
          fill: 'transparent',
          stroke: 'transparent',
        },
      });

      this.addShape(image);
    },
  },
};
</script>

<style lang="scss">
.flowchart-sidebar {
  &__button {
    min-width: 75px !important;
  }
}
</style>
