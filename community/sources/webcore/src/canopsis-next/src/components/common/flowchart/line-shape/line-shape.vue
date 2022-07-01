<template lang="pug">
  g
    slot
    points-path(
      v-bind="shape.properties",
      :points="shape.points",
      :marker-end="markerEnd",
      :marker-start="markerStart",
      pointer-events="none"
    )
    text-editor(
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
  computed: {
    labelPosition() {
      const { points } = this.shape;
      const halfLength = points.length / 2;

      if (points.length % 2 !== 0) {
        return points[Math.floor(halfLength)];
      }

      const p1 = points[halfLength - 1];
      const p2 = points[halfLength];

      return {
        x: (p2.x + p1.x) / 2,
        y: (p2.y + p1.y) / 2,
      };
    },
  },
};
</script>
