<template>
  <g @dblclick.stop="enableEditingMode">
    <ellipse
      v-bind="shape.properties"
      :cx="centerX"
      :cy="centerY"
      :rx="radiusX"
      :ry="radiusY"
      :cursor="readonly ? '' : 'move'"
      v-on="$listeners"
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

export default {
  components: { TextEditor },
  mixins: [flowchartTextEditorMixin],
  props: {
    shape: {
      type: Object,
      required: true,
    },
    selected: {
      type: Boolean,
      default: false,
    },
    readonly: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    radiusX() {
      return this.shape.width / 2;
    },

    radiusY() {
      return this.shape.height / 2;
    },

    centerX() {
      return this.shape.x + this.radiusX;
    },

    centerY() {
      return this.shape.y + this.radiusY;
    },
  },
};
</script>
