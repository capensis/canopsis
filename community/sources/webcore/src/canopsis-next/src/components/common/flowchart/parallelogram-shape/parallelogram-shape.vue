<template lang="pug">
  g
    parallelogram-figure(
      v-bind="shape.properties",
      :width="shape.width",
      :height="shape.height",
      :offset="shape.offset",
      :x="shape.x",
      :y="shape.y"
    )
    text-editor(
      ref="editor",
      v-bind="shape.textProperties",
      :value="shape.text",
      :y="shape.y",
      :x="shape.x",
      :width="shape.width",
      :height="shape.height",
      :editable="editing",
      @blur="disableEditingMode"
    )
    parallelogram-shape-selection(
      v-if="!readonly",
      :selected="selected",
      :connecting="connecting",
      :parallelogram="shape",
      :pointer-events="editing ? 'none' : 'all'",
      @update="$listeners.update",
      @dblclick="enableEditingMode",
      @mousedown="$listeners.mousedown",
      @mouseup="$listeners.mouseup",
      @connected="$listeners.connected",
      @connecting="$listeners.connecting",
      @unconnect="$listeners.unconnect"
    )
</template>

<script>
import { flowchartTextEditorMixin } from '@/mixins/flowchart/text-editor';

import TextEditor from '../common/text-editor.vue';
import ParallelogramFigure from '../common/parallelogram-figure.vue';

import ParallelogramShapeSelection from './parallelogram-shape-selection.vue';

export default {
  components: { ParallelogramFigure, ParallelogramShapeSelection, TextEditor },
  mixins: [flowchartTextEditorMixin],
  props: {
    shape: {
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
};
</script>
