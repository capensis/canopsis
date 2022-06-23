<template lang="pug">
  g
    circle-figure(
      v-bind="circle.style",
      :x="circle.x",
      :y="circle.y",
      :diameter="circle.diameter",
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
      :align-center="circle.alignCenter",
      :justify-center="circle.justifyCenter",
      @blur="disableEditingMode"
    )
    circle-shape-selection(
      v-if="!readonly",
      :selected="selected",
      :connecting="connecting",
      :circle="circle",
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

import CircleShapeSelection from '../circle-shape/circle-shape-selection.vue';
import TextEditor from '../common/text-editor.vue';
import CircleFigure from '../common/circle-figure.vue';

export default {
  components: { CircleFigure, CircleShapeSelection, TextEditor },
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
    onResize(circle) {
      this.updateModel({
        ...this.circle,
        ...circle,
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
