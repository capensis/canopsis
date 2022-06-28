<template lang="pug">
  g
    circle-figure(
      v-bind="shape.properties",
      :x="shape.x",
      :y="shape.y",
      :diameter="shape.diameter",
      pointer-events="all"
    )
    text-editor(
      ref="editor",
      v-bind="shape.textProperties",
      :value="shape.text",
      :y="shape.y",
      :x="shape.x",
      :width="shape.diameter",
      :height="shape.diameter",
      :editable="editing",
      @blur="disableEditingMode"
    )
    circle-shape-selection(
      v-if="!readonly",
      :selected="selected",
      :connecting="connecting",
      :circle="shape",
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
import CircleFigure from '../common/circle-figure.vue';

import CircleShapeSelection from './circle-shape-selection.vue';

export default {
  components: { CircleFigure, CircleShapeSelection, TextEditor },
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
    radius() {
      return this.shape.diameter / 2;
    },

    centerX() {
      return this.shape.x + this.radius;
    },

    centerY() {
      return this.shape.y + this.radius;
    },
  },
};
</script>
