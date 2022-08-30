<template lang="pug">
  g(@dblclick="enableEditingMode", @contextmenu="$listeners.contextmenu")
    document-figure(
      v-bind="shape.properties",
      :width="shape.width",
      :height="shape.height",
      :x="shape.x",
      :y="shape.y",
      :offset="shape.offset",
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
import DocumentFigure from '../common/document-figure.vue';

export default {
  components: { DocumentFigure, TextEditor },
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
