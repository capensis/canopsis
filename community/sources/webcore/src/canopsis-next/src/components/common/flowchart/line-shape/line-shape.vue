<template lang="pug">
  g
    slot
    points-path(
      ref="path",
      v-bind="shape.properties",
      :points="shape.points",
      :type="shape.lineType",
      :marker-end="markerEnd",
      :marker-start="markerStart",
      pointer-events="none"
    )
    text-editor(
      v-if="labelPosition",
      ref="editor",
      v-bind="shape.textProperties",
      :value="shape.text",
      :y="labelPosition.y - shape.textProperties.fontSize / 2",
      :x="labelPosition.x",
      :editable="editing",
      align-center,
      justify-center,
      @blur="disableEditingMode"
    )
    line-shape-selection(
      v-if="!readonly",
      :selected="selected",
      :line="shape",
      @dblclick="enableEditingMode",
      @mousedown="$listeners.mousedown",
      @mouseup="$listeners.mouseup",
      @update="$listeners.update",
      @edit:point="$listeners['edit:point']"
    )
</template>

<script>
import { flowchartTextEditorMixin } from '@/mixins/flowchart/text-editor';

import PointsPath from '../common/points-path.vue';
import TextEditor from '../common/text-editor.vue';

import LineShapeSelection from './line-shape-selection.vue';

export default {
  components: { TextEditor, PointsPath, LineShapeSelection },
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
    cornerOffset: {
      type: Number,
      default: 0,
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
