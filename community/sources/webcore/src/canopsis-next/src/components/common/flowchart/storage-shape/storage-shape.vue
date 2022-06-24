<template lang="pug">
  g
    storage-figure(
      v-bind="shape.style",
      :width="shape.width",
      :height="shape.height",
      :radius="shape.radius",
      :x="shape.x",
      :y="shape.y"
    )
    text-editor(
      ref="editor",
      :value="shape.text",
      :y="shape.y",
      :x="shape.x",
      :width="shape.width",
      :height="shape.height",
      :editable="editing",
      :align-center="shape.alignCenter",
      :justify-center="shape.justifyCenter",
      @blur="disableEditingMode"
    )
    storage-shape-selection(
      v-if="!readonly",
      :selected="selected",
      :connecting="connecting",
      :storage="shape",
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
import StorageFigure from '../common/storage-figure.vue';

import StorageShapeSelection from './storage-shape-selection.vue';

export default {
  components: { StorageFigure, StorageShapeSelection, TextEditor },
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
  computed: {
    diameter() {
      return this.shape.radius * 2;
    },
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
