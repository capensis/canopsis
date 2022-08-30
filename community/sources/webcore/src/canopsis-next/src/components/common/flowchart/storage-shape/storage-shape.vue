<template lang="pug">
  g(@dblclick="enableEditingMode", @contextmenu="$listeners.contextmenu")
    storage-figure(
      v-bind="shape.properties",
      :width="shape.width",
      :height="shape.height",
      :radius="shape.radius",
      :x="shape.x",
      :y="shape.y",
      :cursor="readonly ? '' : 'move'",
      @mousedown.stop="$listeners.mousedown",
      @mouseup="$listeners.mouseup"
    )
    text-editor(
      ref="editor",
      v-bind="shape.textProperties",
      :value="shape.text",
      :y="shape.y",
      :x="shape.x",
      :width="shape.width",
      :height="shape.height",
      :editable="editing",
      @blur="disableEditingMode"
    )
</template>

<script>
import { flowchartTextEditorMixin } from '@/mixins/flowchart/text-editor';

import TextEditor from '../common/text-editor.vue';
import StorageFigure from '../common/storage-figure.vue';

export default {
  components: { StorageFigure, TextEditor },
  mixins: [flowchartTextEditorMixin],
  props: {
    shape: {
      type: Object,
      required: true,
    },
    readonly: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
