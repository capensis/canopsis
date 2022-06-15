<template lang="pug">
  g
    component(is="foreignObject", :y="text.y", :x="text.x", :width="text.width", :height="text.height")
      div(ref="textEditor", :contenteditable="editing") {{ text.text }}
    rect-shape-selection(
      :selected="selected",
      :rect="text",
      @resize="onResize",
      @dblclick="enableEditingMode",
      @mousedown="$emit('mousedown', $event)",
      @mouseup="$emit('mouseup', $event)"
    )
</template>

<script>
import { formBaseMixin } from '@/mixins/form';

import RectShapeSelection from '../rect-shape/rect-shape-selection.vue';

export default {
  components: { RectShapeSelection },
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

      this.$refs.textEditor.focus();
    },
  },
};
</script>
