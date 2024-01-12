<template>
  <g @dblclick.stop="enableEditingMode">
    <process-figure
      v-bind="shape.properties"
      v-on="$listeners"
      :width="shape.width"
      :height="shape.height"
      :x="shape.x"
      :y="shape.y"
      :offset="shape.offset"
      :cursor="readonly ? '' : 'move'"
    />
    <text-editor
      ref="editor"
      v-bind="shape.textProperties"
      :value="shape.text"
      :y="shape.y"
      :x="shape.x"
      :width="shape.width"
      :height="shape.height"
      :editable="editing"
      @blur="disableEditingMode"
    />
  </g>
</template>

<script>
import { flowchartTextEditorMixin } from '@/mixins/flowchart/text-editor';

import TextEditor from '../../common/text-editor.vue';
import ProcessFigure from '../../figures/process-figure.vue';

export default {
  components: { ProcessFigure, TextEditor },
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
