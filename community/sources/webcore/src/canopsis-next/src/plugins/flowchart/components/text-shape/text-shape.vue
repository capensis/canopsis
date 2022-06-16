<template lang="pug">
  g
    text-editor(
      ref="editor",
      :value="text.text",
      :y="text.y",
      :x="text.x",
      :width="text.width",
      :height="text.height",
      :center="text.center",
      :editable="editing",
      @blur="disableEditingMode"
    )
    rect-shape-selection(
      :selected="selected",
      :rect="text",
      :pointer-events="editing ? 'none' : 'all'",
      @resize="onResize",
      @dblclick="enableEditingMode",
      @mousedown="$emit('mousedown', $event)",
      @mouseup="$emit('mouseup', $event)"
    )
</template>

<script>
import { formBaseMixin } from '@/mixins/form';

import RectShapeSelection from '../rect-shape/rect-shape-selection.vue';
import TextEditor from '../common/text-editor.vue';

export default {
  components: { RectShapeSelection, TextEditor },
  mixins: [formBaseMixin],
  model: {
    prop: 'text',
    event: 'input',
  },
  props: {
    text: {
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
      const { x, y } = this.text;

      this.updateModel({
        ...this.text,

        x: (x - oldOffset.x) + newOffset.x,
        y: (y - oldOffset.y) + newOffset.y,
      });
    },

    onResize(rect) {
      this.updateModel({ ...this.text, ...rect });
    },

    enableEditingMode() {
      this.editing = true;

      this.$refs.editor.focus();
    },

    disableEditingMode() {
      this.editing = false;
    },
  },
};
</script>
