<template lang="pug">
  g
    rhombus-figure(
      v-bind="rhombus.style",
      :width="rhombus.width",
      :height="rhombus.height",
      :x="rhombus.x",
      :y="rhombus.y"
    )
    text-editor(
      ref="editor",
      :value="rhombus.text",
      :y="rhombus.y",
      :x="rhombus.x",
      :width="rhombus.width",
      :height="rhombus.height",
      :editable="editing",
      :align-center="rhombus.alignCenter",
      :justify-center="rhombus.justifyCenter",
      @blur="disableEditingMode"
    )
    rhombus-shape-selection(
      v-if="!readonly",
      :selected="selected",
      :connecting="connecting",
      :rhombus="rhombus",
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
import RhombusFigure from '../common/rhombus-figure.vue';

import RhombusShapeSelection from './rhombus-shape-selection.vue';

export default {
  components: { RhombusFigure, RhombusShapeSelection, TextEditor },
  mixins: [formBaseMixin],
  model: {
    prop: 'rhombus',
    event: 'input',
  },
  props: {
    rhombus: {
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
      const { x, y } = this.rhombus;

      this.updateModel({
        ...this.rhombus,

        x: (x - oldOffset.x) + newOffset.x,
        y: (y - oldOffset.y) + newOffset.y,
      });
    },

    onResize(rhombus) {
      this.updateModel({ ...this.rhombus, ...rhombus });
    },

    enableEditingMode() {
      this.editing = true;

      this.$refs.editor.focus();
    },

    disableEditingMode(event) {
      this.updateModel({
        ...this.rhombus,
        text: event.target.innerHTML,
      });

      this.editing = false;
    },
  },
};
</script>
