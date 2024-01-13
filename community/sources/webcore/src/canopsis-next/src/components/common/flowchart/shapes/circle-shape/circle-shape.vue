<template>
  <g @dblclick.stop="enableEditingMode">
    <circle-figure
      v-bind="shape.properties"
      v-on="$listeners"
      :x="shape.x"
      :y="shape.y"
      :diameter="shape.diameter"
      :cursor="readonly ? '' : 'move'"
    />
    <text-editor
      ref="editor"
      v-bind="shape.textProperties"
      :value="shape.text"
      :y="shape.y"
      :x="shape.x"
      :width="shape.diameter"
      :height="shape.diameter"
      :editable="editing"
      @blur="disableEditingMode"
    />
  </g>
</template>

<script>
import { flowchartTextEditorMixin } from '@/mixins/flowchart/text-editor';

import TextEditor from '../../common/text-editor.vue';
import CircleFigure from '../../figures/circle-figure.vue';

export default {
  components: { CircleFigure, TextEditor },
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
