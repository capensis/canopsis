<template lang="pug">
  g
    ellipse(
      v-bind="shape.properties",
      :cx="centerX",
      :cy="centerY",
      :rx="radiusX",
      :ry="radiusY"
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
    rect-shape-selection(
      v-if="!readonly",
      :selected="selected",
      :connecting="connecting",
      :rect="shape",
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

import RectShapeSelection from '../rect-shape/rect-shape-selection.vue';
import TextEditor from '../common/text-editor.vue';

export default {
  components: { RectShapeSelection, TextEditor },
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
  computed: {
    radiusX() {
      return this.shape.width / 2;
    },

    radiusY() {
      return this.shape.height / 2;
    },

    centerX() {
      return this.shape.x + this.radiusX;
    },

    centerY() {
      return this.shape.y + this.radiusY;
    },
  },
};
</script>
