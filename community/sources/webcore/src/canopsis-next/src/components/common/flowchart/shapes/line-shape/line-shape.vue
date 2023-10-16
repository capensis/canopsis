<template>
  <g @dblclick.stop="enableEditingMode">
    <slot />
    <points-line-path
      ref="path"
      v-bind="shape.properties"
      :points="shape.points"
      :type="shape.lineType"
      :marker-end="markerEnd"
      :marker-start="markerStart"
    />
    <points-line-path
      v-on="$listeners"
      :points="shape.points"
      :type="shape.lineType"
      :cursor="readonly ? '' : 'move'"
      stroke-width="10"
    />
    <text-editor
      v-if="labelPosition"
      ref="editor"
      v-bind="shape.textProperties"
      :value="shape.text"
      :y="labelPosition.y - shape.textProperties.fontSize / 2"
      :x="labelPosition.x"
      :editable="editing"
      align-center
      justify-center
      @blur="disableEditingMode"
    />
  </g>
</template>

<script>
import { flowchartTextEditorMixin } from '@/mixins/flowchart/text-editor';

import PointsLinePath from '../../common/points-line-path.vue';
import TextEditor from '../../common/text-editor.vue';

export default {
  components: { TextEditor, PointsLinePath },
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
    markerStart: {
      type: String,
      required: false,
    },
    markerEnd: {
      type: String,
      required: false,
    },
  },
  data() {
    return {
      labelPosition: null,
    };
  },
  watch: {
    'shape.points': {
      deep: true,
      handler() {
        this.$nextTick(this.calculateLabelPosition);
      },
    },
  },
  mounted() {
    this.calculateLabelPosition();
  },
  methods: {
    calculateLabelPosition() {
      this.labelPosition = this.$refs.path.getCenterPoint();
    },
  },
};
</script>
