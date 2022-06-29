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

      this.$set(this.shapes, rect._id, rect);
    },

    addRoundedRectangle() {
      const rect = generateRectShape({
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

      this.$set(this.shapes, rect._id, rect);
    },

    addLine() {
      const line = generateLineShape({
        text: 'Line',
        properties: {
          stroke: 'black',
          'stroke-width': 2,
        },
      });

      this.$set(this.shapes, line._id, line);
    },

    addDashedLine() {
      const line = generateLineShape({
        text: 'Dashed line',
        properties: {
          stroke: 'black',
          'stroke-width': 2,
        },
      });

      this.$set(this.shapes, line._id, line);
    },

    addDashedArrowLine() {
      const arrowLine = generateArrowLineShape({
        text: 'Dashed arrow line',
        properties: {
          stroke: 'black',
          'stroke-width': 2,
          'stroke-dasharray': '4 4',
        },
      });

      this.$set(this.shapes, arrowLine._id, arrowLine);
    },

    addArrowLine() {
      const arrowLine = generateArrowLineShape({
        text: 'Arrow line',
        properties: {
          stroke: 'black',
          'stroke-width': 2,
        },
      });

      this.$set(this.shapes, arrowLine._id, arrowLine);
    },

    addBidirectionalArrowLine() {
      const bidirectionalArrowLine = generateBidirectionalArrowLineShape({
        text: 'Bidirectional arrow line',
        properties: {
          stroke: 'black',
          'stroke-width': 2,
        },
      });

      this.$set(this.shapes, bidirectionalArrowLine._id, bidirectionalArrowLine);
    },

    addDashedBidirectionalArrowLine() {
      const bidirectionalArrowLine = generateBidirectionalArrowLineShape({
        text: 'Dashed bidirectional arrow line',
        properties: {
          stroke: 'black',
          'stroke-width': 2,
          'stroke-dasharray': '4 4',
        },
      });

      this.$set(this.shapes, bidirectionalArrowLine._id, bidirectionalArrowLine);
    },

    addCircle() {
      const circle = generateCircleShape({
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

      this.$set(this.shapes, circle._id, circle);
    },

    addEllipse() {
      const ellipse = generateEllipseShape({
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

      this.$set(this.shapes, ellipse._id, ellipse);
    },

    addRhombus() {
      const rhombus = generateRhombusShape({
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

      this.$set(this.shapes, rhombus._id, rhombus);
    },

    addParallelogram() {
      const parallelogram = generateParallelogramShape({
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

      this.$set(this.shapes, parallelogram._id, parallelogram);
    },

    addSquare() {
      const square = generateRectShape({
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

      this.$set(this.shapes, square._id, square);
    },

    addText() {
      const text = generateRectShape({
        text: 'Text',
        textProperties: {
          alignCenter: true,
          justifyCenter: true,
        },
        properties: {
          fill: 'transparent',
        },
      });

      this.$set(this.shapes, text._id, text);
    },

    addTextbox() {
      const textbox = generateRectShape({
        text: '<h2>Heading</h2><p>Paragraph</p>',
        textProperties: {
          hidden: true,
        },
        properties: {
          fill: 'transparent',
        },
      });

      this.$set(this.shapes, textbox._id, textbox);
    },

    addStorage() {
      const storage = generateStorageShape({
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

      this.$set(this.shapes, storage._id, storage);
    },

    async addImage([file]) {
      const src = await getFileDataUrlContent(file);

      const img = document.createElement('img');

      img.onload = () => {
        const image = generateImageShape({
          text: file.name,
          aspectRatio: true,
          x: 0,
          y: 0,
          width: img.width ?? 100,
          height: img.height ?? 100,
          src,
        });

        this.$set(this.shapes, image._id, image);

        this.$refs.fileSelector.clear();
      };

      img.src = src;
    },
  },
};
</script>

<style lang="scss">
.flowchart {
  display: flex;

  &__editor {
    flex-grow: 1;
  }
}
</style>
