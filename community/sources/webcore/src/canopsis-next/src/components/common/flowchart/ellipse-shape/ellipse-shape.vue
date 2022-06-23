<template lang="pug">
  g
    ellipse(
      v-bind="ellipse.style",
      :cx="centerX",
      :cy="centerY",
      :rx="radiusX",
      :ry="radiusY"
    )
    text-editor(
      ref="editor",
      :value="ellipse.text",
      :y="ellipse.y",
      :x="ellipse.x",
      :width="ellipse.width",
      :height="ellipse.height",
      :editable="editing",
      :align-center="ellipse.alignCenter",
      :justify-center="ellipse.justifyCenter",
      @blur="disableEditingMode"
    )
    rect-shape-selection(
      v-if="!readonly",
      :selected="selected",
      :connecting="connecting",
      :rect="ellipse",
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

import RectShapeSelection from '../rect-shape/rect-shape-selection.vue';
import TextEditor from '../common/text-editor.vue';

export default {
  components: { RectShapeSelection, TextEditor },
  mixins: [formBaseMixin],
  model: {
    prop: 'ellipse',
    event: 'input',
  },
  props: {
    ellipse: {
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
    radiusX() {
      return this.ellipse.width / 2;
    },

    radiusY() {
      return this.ellipse.height / 2;
    },

    centerX() {
      return this.ellipse.x + this.radiusX;
    },

    centerY() {
      return this.ellipse.y + this.radiusY;
    },
  },
  methods: {
    onResize(ellipse) {
      this.updateModel({ ...this.ellipse, ...ellipse });
    },

    enableEditingMode() {
      this.editing = true;

      this.$refs.editor.focus();
    },

    disableEditingMode(event) {
      this.updateModel({
        ...this.ellipse,
        text: event.target.innerHTML,
      });

      this.editing = false;
    },
  },
};
</script>
