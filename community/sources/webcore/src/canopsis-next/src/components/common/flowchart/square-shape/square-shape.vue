<template lang="pug">
  g
    rect(
      v-bind="shape.style",
      :x="shape.x",
      :y="shape.y",
      :width="shape.size",
      :height="shape.size"
    )
    text-editor(
      ref="editor",
      :value="shape.text",
      :y="shape.y",
      :x="shape.x",
      :width="shape.size",
      :height="shape.size",
      :editable="editing",
      align-center,
      justify-center,
      @blur="disableEditingMode"
    )
    square-shape-selection(
      v-if="!readonly",
      :selected="selected",
      :connecting="connecting",
      :square="shape",
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
import TextEditor from '../common/text-editor.vue';

import SquareShapeSelection from './square-shape-selection.vue';

export default {
  components: { SquareShapeSelection, TextEditor },
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
  data() {
    return {
      editing: false,
    };
  },
  methods: {
    enableEditingMode() {
      this.editing = true;

      this.$refs.editor.focus();
    },

    disableEditingMode(event) {
      this.$emit('update', {
        text: event.target.innerHTML,
      });

      this.editing = false;
    },
  },
};
</script>
