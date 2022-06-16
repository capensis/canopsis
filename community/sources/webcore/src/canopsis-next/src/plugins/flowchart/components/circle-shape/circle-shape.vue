<template lang="pug">
  g
    circle(
      v-bind="circle.style",
      :cx="centerX",
      :cy="centerY",
      :r="radius",
      pointer-events="all"
    )
    text-editor(
      ref="editor",
      :value="circle.text",
      :y="circle.y",
      :x="circle.x",
      :width="circle.diameter",
      :height="circle.diameter",
      :editable="editing",
      center,
      @blur="disableEditingMode"
    )
    circle-shape-selection(
      :selected="selected",
      :circle="circle",
      :pointer-events="editing ? 'none' : 'all'",
      @resize="onResize",
      @dblclick="enableEditingMode",
      @mousedown="$emit('mousedown', $event)",
      @mouseup="$emit('mouseup', $event)"
    )
</template>

<script>
import { formBaseMixin } from '@/mixins/form';

import CircleShapeSelection from './circle-shape-selection.vue';
import TextEditor from '../common/text-editor.vue';

export default {
  components: { CircleShapeSelection, TextEditor },
  mixins: [formBaseMixin],
  model: {
    prop: 'circle',
    event: 'input',
  },
  props: {
    circle: {
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
  computed: {
    radius() {
      return this.circle.diameter / 2;
    },

    centerX() {
      return this.circle.x + this.radius;
    },

    centerY() {
      return this.circle.y + this.radius;
    },
  },
  methods: {
    move(newOffset, oldOffset) {
      const { x, y } = this.circle;

      this.updateModel({
        ...this.circle,

        x: (x - oldOffset.x) + newOffset.x,
        y: (y - oldOffset.y) + newOffset.y,
      });
    },

    onResize(circle) {
      this.updateModel({
        ...this.circle,
        x: circle.x,
        y: circle.y,
        diameter: circle.size,
      });
    },

    enableEditingMode() {
      this.editing = true;

      this.$refs.editor.focus();
    },

    disableEditingMode(event) {
      this.updateModel({
        ...this.circle,
        text: event.target.innerHTML,
      });

      this.editing = false;
    },
  },
};
</script>
