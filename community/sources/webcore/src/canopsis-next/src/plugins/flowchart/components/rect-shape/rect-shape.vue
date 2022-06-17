<template lang="pug">
  g
    rect(
      v-bind="rect.style",
      :x="rect.x",
      :y="rect.y",
      :width="rect.width",
      :height="rect.height"
    )
    text-editor(
      ref="editor",
      :value="rect.text",
      :y="rect.y",
      :x="rect.x",
      :width="rect.width",
      :height="rect.height",
      :editable="editing",
      :center="rect.textCenter",
      @blur="disableEditingMode"
    )
    rect-shape-selection(
      v-if="!readonly",
      :selected="selected",
      :rect="rect",
      :pointer-events="editing ? 'none' : 'all'",
      @resize="onResize",
      @dblclick="enableEditingMode",
      @mousedown="$emit('mousedown', $event)",
      @mouseup="$emit('mouseup', $event)"
    )
</template>

<script>
import { formBaseMixin } from '@/mixins/form';

import RectShapeSelection from './rect-shape-selection.vue';
import TextEditor from '../common/text-editor.vue';

export default {
  components: { RectShapeSelection, TextEditor },
  mixins: [formBaseMixin],
  model: {
    prop: 'rect',
    event: 'input',
  },
  props: {
    rect: {
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
  },
  data() {
    return {
      editing: false,
    };
  },
  methods: {
    move(newOffset, oldOffset) {
      const { x, y } = this.rect;

      this.updateModel({
        ...this.rect,

        x: (x - oldOffset.x) + newOffset.x,
        y: (y - oldOffset.y) + newOffset.y,
      });
    },

    onResize(rect) {
      this.updateModel({ ...this.rect, ...rect });
    },

    enableEditingMode() {
      this.editing = true;

      this.$refs.editor.focus();
    },

    disableEditingMode(event) {
      this.updateModel({
        ...this.rect,
        text: event.target.innerHTML,
      });

      this.editing = false;
    },
  },
};
</script>
