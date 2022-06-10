<template lang="pug">
  g
    foreign-object(width="100%", height="100%")
      div(
        :contenteditable="editable",
        :style="`position: absolute; top: ${text.y}px;left: ${text.x}px;`"
      ) {{ text.text }}
    rect-shape-selection(
      :selected="selected",
      :rect="text",
      @resize="onResize",
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
      editable: false,
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
  },
};
</script>
