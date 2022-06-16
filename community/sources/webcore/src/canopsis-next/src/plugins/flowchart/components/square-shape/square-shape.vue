<template lang="pug">
  g
    rect(
      v-bind="square.style",
      :x="square.x",
      :y="square.y",
      :width="square.size",
      :height="square.size"
    )
    text-editor(
      ref="editor",
      :value="square.text",
      :y="square.y",
      :x="square.x",
      :width="square.size",
      :height="square.size",
      :editable="editing",
      center,
      @blur="disableEditingMode"
    )
    square-shape-selection(
      :selected="selected",
      :square="square",
      :pointer-events="editing ? 'none' : 'all'",
      @resize="onResize",
      @dblclick="enableEditingMode",
      @mousedown="$emit('mousedown', $event)",
      @mouseup="$emit('mouseup', $event)"
    )
</template>

<script>
import { formBaseMixin } from '@/mixins/form';

import SquareShapeSelection from './square-shape-selection.vue';
import TextEditor from '../common/text-editor.vue';

export default {
  components: { SquareShapeSelection, TextEditor },
  mixins: [formBaseMixin],
  model: {
    prop: 'square',
    event: 'input',
  },
  props: {
    square: {
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
  },
  data() {
    return {
      editing: false,
    };
  },
  methods: {
    move(newOffset, oldOffset) {
      const { x, y } = this.square;

      this.updateModel({
        ...this.square,

        x: (x - oldOffset.x) + newOffset.x,
        y: (y - oldOffset.y) + newOffset.y,
      });
    },

    onResize(square) {
      this.updateModel({
        ...this.square,
        x: square.x,
        y: square.y,
        size: square.size,
      });
    },

    enableEditingMode() {
      this.editing = true;

      this.$refs.editor.focus();
    },

    disableEditingMode(event) {
      this.updateModel({
        ...this.square,
        text: event.target.innerHTML,
      });

      this.editing = false;
    },
  },
};
</script>
