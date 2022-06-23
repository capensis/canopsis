<template lang="pug">
  g
    parallelogram-figure(
      v-bind="parallelogram.style",
      :width="parallelogram.width",
      :height="parallelogram.height",
      :offset="parallelogram.offset",
      :x="parallelogram.x",
      :y="parallelogram.y"
    )
    text-editor(
      ref="editor",
      :value="parallelogram.text",
      :y="parallelogram.y",
      :x="parallelogram.x",
      :width="parallelogram.width",
      :height="parallelogram.height",
      :editable="editing",
      :align-center="parallelogram.alignCenter",
      :justify-center="parallelogram.justifyCenter",
      @blur="disableEditingMode"
    )
    parallelogram-shape-selection(
      v-if="!readonly",
      :selected="selected",
      :connecting="connecting",
      :parallelogram="parallelogram",
      :pointer-events="editing ? 'none' : 'all'",
      @resize="onResize",
      @dblclick="enableEditingMode",
      @mousedown="$listeners.mousedown",
      @mouseup="$listeners.mouseup",
      @connected="$listeners.connected",
      @connecting="$listeners.connecting",
      @unconnect="$listeners.unconnect"
    )
</template>

<script>
import { formBaseMixin } from '@/mixins/form';

import TextEditor from '../common/text-editor.vue';
import ParallelogramFigure from '../common/parallelogram-figure.vue';

import ParallelogramShapeSelection from './parallelogram-shape-selection.vue';

export default {
  components: { ParallelogramFigure, ParallelogramShapeSelection, TextEditor },
  mixins: [formBaseMixin],
  model: {
    prop: 'parallelogram',
    event: 'input',
  },
  props: {
    parallelogram: {
      type: Object,
      required: true,
    },
    selected: {
      type: Boolean,
      default: false,
    },
    cornerOffset: {
      type: Number,
      default: 0,
    },
    readonly: {
      type: Boolean,
      default: false,
    },
    connecting: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      editing: false,
    };
  },
  methods: {
    move(newOffset, oldOffset) {
      const { x, y } = this.parallelogram;

      this.updateModel({
        ...this.parallelogram,

        x: (x - oldOffset.x) + newOffset.x,
        y: (y - oldOffset.y) + newOffset.y,
      });
    },

    onResize(parallelogram) {
      this.updateModel({ ...this.parallelogram, ...parallelogram });
    },

    enableEditingMode() {
      this.editing = true;

      this.$refs.editor.focus();
    },

    disableEditingMode(event) {
      this.updateModel({
        ...this.parallelogram,
        text: event.target.innerHTML,
      });

      this.editing = false;
    },
  },
};
</script>
