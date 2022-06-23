<template lang="pug">
  g
    storage-figure(
      v-bind="storage.style",
      :width="storage.width",
      :height="storage.height",
      :radius="storage.radius",
      :x="storage.x",
      :y="storage.y"
    )
    text-editor(
      ref="editor",
      :value="storage.text",
      :y="storage.y",
      :x="storage.x",
      :width="storage.width",
      :height="storage.height",
      :editable="editing",
      :align-center="storage.alignCenter",
      :justify-center="storage.justifyCenter",
      @blur="disableEditingMode"
    )
    storage-shape-selection(
      v-if="!readonly",
      :selected="selected",
      :connecting="connecting",
      :storage="storage",
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

import TextEditor from '../common/text-editor.vue';
import StorageFigure from '../common/storage-figure.vue';

import StorageShapeSelection from '../storage-shape/storage-shape-selection.vue';

export default {
  components: { StorageFigure, StorageShapeSelection, TextEditor },
  mixins: [formBaseMixin],
  model: {
    prop: 'storage',
    event: 'input',
  },
  props: {
    storage: {
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
      return this.storage.radius * 2;
    },
  },
  methods: {
    onResize(storage) {
      this.updateModel({ ...this.storage, ...storage });
    },

    enableEditingMode() {
      this.editing = true;

      this.$refs.editor.focus();
    },

    disableEditingMode(event) {
      this.updateModel({
        ...this.storage,
        text: event.target.innerHTML,
      });

      this.editing = false;
    },
  },
};
</script>
