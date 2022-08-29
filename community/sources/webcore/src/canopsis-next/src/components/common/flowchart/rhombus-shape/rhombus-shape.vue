<template lang="pug">
  g(@contextmenu="$listeners.contextmenu")
    rhombus-figure(
      v-bind="shape.properties",
      :width="shape.width",
      :height="shape.height",
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
    rhombus-shape-selection(
      v-if="!readonly",
      :selected="selected",
      :connecting="connecting",
      :rhombus="shape",
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
import RhombusFigure from '../common/rhombus-figure.vue';

import RhombusShapeSelection from './rhombus-shape-selection.vue';

export default {
  components: { RhombusFigure, RhombusShapeSelection, TextEditor },
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
