<template>
  <v-navigation-drawer
    class="flowchart-sidebar"
    width="100%"
    permanent
    touchless
    absolute
  >
    <c-layout column>
      <slot name="prepend" />
      <flowchart-color-field
        v-model:value="localBackgroundColor"
        :label="$t('flowchart.backgroundColor')"
        :palette="backgroundColors"
        class="flowchart-sidebar__additional-item px-6"
        hide-checkbox
      />
    </c-layout>
    <v-divider />
    <v-expansion-panels variant="accordion" multiple>
      <v-expansion-panel>
        <v-expansion-panel-title>
          <span class="v-label">
            {{ $t('flowchart.shape', 2) }}
          </span>
        </v-expansion-panel-title>
        <v-divider />
        <v-expansion-panel-text>
          <c-layout wrap>
            <v-tooltip
              v-for="(button, index) in buttons"
              :key="index"
              z-index="230"
              location="top"
            >
              <template #activator="{ props: activatorProps }">
                <v-btn
                  v-bind="activatorProps"
                  class="ma-0 pa-0 flowchart-sidebar__button"
                  variant="text"
                  size="large"
                  @click="button.action"
                >
                  <component :is="button.icon" class="pa-1 text-disabled" />
                </v-btn>
              </template>
              <span>{{ button.label }}</span>
            </v-tooltip>
            <file-selector ref="fileSelector" hide-details @change="addImage">
              <template #activator="{ props: fileActivatorProps }">
                <v-tooltip z-index="230" location="top">
                  <template #activator="{ props: tooltipActivatorProps }">
                    <v-btn
                      v-bind="{ ...fileActivatorProps, ...tooltipActivatorProps }"
                      class="ma-0 pa-0 flowchart-sidebar__button"
                      variant="text"
                      size="large"
                    >
                      <image-shape-icon class="text-grey-darken-3" />
                    </v-btn>
                  </template>
                  <span>{{ $t('flowchart.shapes.image') }}</span>
                </v-tooltip>
              </template>
            </file-selector>
          </c-layout>
          <v-divider />
        </v-expansion-panel-text>
      </v-expansion-panel>
      <v-expansion-panel>
        <v-expansion-panel-title>
          <span class="v-label">
            {{ $t('flowchart.icons') }}
          </span>
        </v-expansion-panel-title>
        <v-divider />
        <v-expansion-panel-text>
          <c-layout v-for="group in iconGroups" :key="group.name" wrap>
            <v-btn
              v-for="(icon, index) in group.icons"
              :key="index"
              class="flowchart-sidebar__button-icon ma-0"
              variant="text"
              @click="addIconAsset(icon)"
            >
              <span v-html="icon" class="text-disabled flowchart-sidebar__button-svg" />
            </v-btn>
            <v-col cols="12">
              <v-divider />
            </v-col>
          </c-layout>
          <v-divider />
        </v-expansion-panel-text>
      </v-expansion-panel>
    </v-expansion-panels>
  </v-navigation-drawer>
</template>

<script>
import { COLORS } from '@/config';
import { LINE_TYPES } from '@/constants';

import { getFileDataUrlContent } from '@/helpers/file/file-select';
import { getImageProperties } from '@/helpers/file/image';
import { generatePoint } from '@/helpers/flowchart/points';
import {
  arrowLineShapeToForm,
  bidirectionalArrowLineShapeToForm,
  circleShapeToForm,
  ellipseShapeToForm,
  imageShapeToForm,
  lineShapeToForm,
  parallelogramShapeToForm,
  processShapeToForm,
  documentShapeToForm,
  rectShapeToForm,
  rhombusShapeToForm,
  storageShapeToForm,
} from '@/helpers/flowchart/shapes';

import { formBaseMixin } from '@/mixins/form';

import FileSelector from '@/components/forms/fields/file-selector.vue';

import assetGroups from './assets';
import FlowchartColorField from './fields/flowchart-color-field.vue';
import ArrowLineShapeIcon from './icons/arrow-line-shape.vue';
import BidirectionalArrowLineShapeIcon from './icons/bidirectional-arrow-line-shape.vue';
import BidirectionalCurveArrowLineShape from './icons/bidirectional-curve-arrow-line-shape.vue';
import CircleShapeIcon from './icons/circle-shape.vue';
import CurveArrowLineShapeIcon from './icons/curve-arrow-line-shape.vue';
import CurveLineShapeIcon from './icons/curve-line-shape.vue';
import DocumentShapeIcon from './icons/document-shape.vue';
import EllipseShapeIcon from './icons/ellipse-shape.vue';
import ImageShapeIcon from './icons/image-shape.vue';
import LineShapeIcon from './icons/line-shape.vue';
import ParallelogramShapeIcon from './icons/parallelogram-shape.vue';
import ProcessShapeIcon from './icons/process-shape.vue';
import RectShapeIcon from './icons/rect-shape.vue';
import RhombusShapeIcon from './icons/rhombus-shape.vue';
import RoundedRectShapeIcon from './icons/rounded-rect-shape.vue';
import SquareShapeIcon from './icons/square-shape.vue';
import StorageShapeIcon from './icons/storage-shape.vue';
import TextShapeIcon from './icons/text-shape.vue';
import TextboxShapeIcon from './icons/textbox-shape.vue';

export default {
  components: { FlowchartColorField, FileSelector, ImageShapeIcon },
  mixins: [formBaseMixin],
  model: {
    prop: 'shapes',
    event: 'input',
  },
  emits: ['input', 'update:shapes', 'update:selected', 'update:backgroundColor'],
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
    backgroundColor: {
      type: String,
      required: false,
    },
    pointDistance: {
      type: Number,
      default: 150,
    },
  },
  computed: {
    localBackgroundColor: {
      get() {
        return this.backgroundColor;
      },

      set(value) {
        this.$emit('update:backgroundColor', value);
      },
    },

    buttons() {
      return [
        {
          icon: RectShapeIcon,
          label: this.$t('flowchart.shapes.rectangle'),
          action: this.addRectangle,
        },
        {
          icon: RoundedRectShapeIcon,
          label: this.$t('flowchart.shapes.roundedRectangle'),
          action: this.addRoundedRectangle,
        },
        {
          icon: SquareShapeIcon,
          label: this.$t('flowchart.shapes.square'),
          action: this.addSquare,
        },
        {
          icon: RhombusShapeIcon,
          label: this.$t('flowchart.shapes.rhombus'),
          action: this.addRhombus,
        },
        {
          icon: CircleShapeIcon,
          label: this.$t('flowchart.shapes.circle'),
          action: this.addCircle,
        },
        {
          icon: EllipseShapeIcon,
          label: this.$t('flowchart.shapes.ellipse'),
          action: this.addEllipse,
        },
        {
          icon: ParallelogramShapeIcon,
          label: this.$t('flowchart.shapes.parallelogram'),
          action: this.addParallelogram,
        },
        {
          icon: ProcessShapeIcon,
          label: this.$t('flowchart.shapes.process'),
          action: this.addProcess,
        },
        {
          icon: DocumentShapeIcon,
          label: this.$t('flowchart.shapes.document'),
          action: this.addDocument,
        },
        {
          icon: StorageShapeIcon,
          label: this.$t('flowchart.shapes.storage'),
          action: this.addStorage,
        },
        {
          icon: CurveLineShapeIcon,
          label: this.$t('flowchart.shapes.curve'),
          action: this.addCurveLine,
        },
        {
          icon: CurveArrowLineShapeIcon,
          label: this.$t('flowchart.shapes.curveArrow'),
          action: this.addCurveArrowLine,
        },
        {
          icon: BidirectionalCurveArrowLineShape,
          label: this.$t('flowchart.shapes.bidirectionalCurve'),
          action: this.addBidirectionalCurveArrowLine,
        },
        {
          icon: LineShapeIcon,
          label: this.$t('flowchart.shapes.line'),
          action: this.addLine,
        },
        {
          icon: ArrowLineShapeIcon,
          label: this.$t('flowchart.shapes.arrowLine'),
          action: this.addArrowLine,
        },
        {
          icon: BidirectionalArrowLineShapeIcon,
          label: this.$t('flowchart.shapes.bidirectionalArrowLine'),
          action: this.addBidirectionalArrowLine,
        },
        {
          icon: TextShapeIcon,
          label: this.$t('flowchart.shapes.text'),
          action: this.addText,
        },
        {
          icon: TextboxShapeIcon,
          label: this.$t('flowchart.shapes.textbox'),
          action: this.addTextbox,
        },
      ];
    },

    halfPointDistance() {
      return this.pointDistance / 2;
    },

    iconGroups() {
      return Object.entries(assetGroups).map(([name, icons]) => ({
        name,
        icons,
      }));
    },

    viewBoxCenter() {
      return {
        x: this.viewBox.x + this.viewBox.width / 2,
        y: this.viewBox.y + this.viewBox.height / 2,
      };
    },

    centerRectProperties() {
      return {
        x: this.viewBoxCenter.x - this.halfPointDistance,
        y: this.viewBoxCenter.y - this.halfPointDistance,
        width: this.pointDistance,
        height: this.pointDistance,
      };
    },

    centerCircleProperties() {
      return {
        x: this.viewBoxCenter.x - this.halfPointDistance,
        y: this.viewBoxCenter.y - this.halfPointDistance,
        diameter: this.pointDistance,
      };
    },

    centerLinePoints() {
      return [
        generatePoint({
          x: this.viewBoxCenter.x - this.halfPointDistance,
          y: this.viewBoxCenter.y + this.halfPointDistance,
        }),
        generatePoint({
          x: this.viewBoxCenter.x + this.halfPointDistance,
          y: this.viewBoxCenter.y - this.halfPointDistance,
        }),
      ];
    },

    backgroundColors() {
      return COLORS.flowchart.background;
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
      const rect = rectShapeToForm({
        ...this.centerRectProperties,
        text: this.$t('flowchart.shapes.rectangle'),
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
      const rect = rectShapeToForm({
        ...this.centerRectProperties,
        text: this.$t('flowchart.shapes.roundedRectangle'),
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
      const line = lineShapeToForm({
        points: this.centerLinePoints,
        properties: {
          stroke: 'black',
          'stroke-width': 2,
        },
      });

      this.addShape(line);
    },

    addCurveLine() {
      const line = lineShapeToForm({
        points: this.centerLinePoints,
        lineType: LINE_TYPES.horizontalCurve,
        properties: {
          stroke: 'black',
          'stroke-width': 2,
        },
      });

      this.addShape(line);
    },

    addCurveArrowLine() {
      const arrowLine = arrowLineShapeToForm({
        points: this.centerLinePoints,
        lineType: LINE_TYPES.horizontalCurve,
        properties: {
          stroke: 'black',
          'stroke-width': 2,
        },
      });

      this.addShape(arrowLine);
    },

    addBidirectionalCurveArrowLine() {
      const bidirectionalArrowLine = bidirectionalArrowLineShapeToForm({
        points: this.centerLinePoints,
        lineType: LINE_TYPES.horizontalCurve,
        properties: {
          stroke: 'black',
          'stroke-width': 2,
        },
      });

      this.addShape(bidirectionalArrowLine);
    },

    addArrowLine() {
      const arrowLine = arrowLineShapeToForm({
        points: this.centerLinePoints,
        properties: {
          stroke: 'black',
          'stroke-width': 2,
        },
      });

      this.addShape(arrowLine);
    },

    addBidirectionalArrowLine() {
      const bidirectionalArrowLine = bidirectionalArrowLineShapeToForm({
        points: this.centerLinePoints,
        properties: {
          stroke: 'black',
          'stroke-width': 2,
        },
      });

      this.addShape(bidirectionalArrowLine);
    },

    addCircle() {
      const circle = circleShapeToForm({
        ...this.centerCircleProperties,
        text: this.$t('flowchart.shapes.circle'),
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
      const ellipse = ellipseShapeToForm({
        ...this.centerRectProperties,
        text: this.$t('flowchart.shapes.ellipse'),
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
      const rhombus = rhombusShapeToForm({
        ...this.centerRectProperties,
        text: this.$t('flowchart.shapes.rhombus'),
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
      const parallelogram = parallelogramShapeToForm({
        ...this.centerRectProperties,
        text: this.$t('flowchart.shapes.parallelogram'),
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
      const parallelogram = processShapeToForm({
        ...this.centerRectProperties,
        text: this.$t('flowchart.shapes.process'),
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
      const parallelogram = documentShapeToForm({
        ...this.centerRectProperties,
        text: this.$t('flowchart.shapes.document'),
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
      const square = rectShapeToForm({
        ...this.centerRectProperties,
        text: this.$t('flowchart.shapes.square'),
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
      const text = rectShapeToForm({
        ...this.centerRectProperties,
        text: this.$t('flowchart.shapes.text'),
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
      const textbox = rectShapeToForm({
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
      const storage = storageShapeToForm({
        ...this.centerRectProperties,
        text: this.$t('flowchart.shapes.storage'),
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

      this.addShape(storage);
    },

    async addImage([file]) {
      const src = await getFileDataUrlContent(file);
      const { width, height } = await getImageProperties(src);

      const maxImageWidth = this.viewBox.width * 0.75;
      const maxImageHeight = this.viewBox.height * 0.75;

      let imageWidth = width;
      let imageHeight = height;

      if (width > maxImageWidth) {
        imageWidth = maxImageWidth;
        imageHeight = height / (width / imageWidth);
      }

      if (imageHeight > maxImageHeight) {
        imageHeight = maxImageHeight;
        imageWidth = width / (height / imageHeight);
      }

      const image = imageShapeToForm({
        ...this.centerRectProperties,
        x: this.viewBoxCenter.x - imageWidth / 2,
        y: this.viewBoxCenter.y - imageHeight / 2,
        width: imageWidth,
        height: imageHeight,
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

    async addIconAsset(svg) {
      const image = imageShapeToForm({
        ...this.centerRectProperties,
        svg,
        aspectRatio: true,
        properties: {
          fill: 'black',
        },
      });

      this.addShape(image);
    },
  },
};
</script>

<style lang="scss">
.flowchart-sidebar {
  .v-navigation-drawer__border {
    z-index: 1;
  }

  &__button {
    min-width: 75px !important;
  }

  &__button-icon {
    min-width: unset !important;
    width: 50px !important;
    height: 50px !important;
    padding: 0 !important;

    svg {
      width: 30px;
      height: 30px;
    }
  }

  &__additional-item {
    height: 51px;
  }
}
</style>
