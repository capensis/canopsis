<template lang="pug">
  g(@dblclick.stop="enableEditingMode")
    rect(
      v-bind="shape.properties",
      :x="shape.x",
      :y="shape.y",
      :width="shape.width",
      :height="shape.height"
    )
    image(
      v-on="$listeners",
      :href="shape.src",
      :x="shape.x",
      :y="shape.y",
      :width="shape.width",
      :height="shape.height",
      :cursor="readonly ? '' : 'move'"
    )
    text-editor(
      ref="editor",
      v-bind="shape.textProperties",
      :value="shape.text",
      :y="shape.y + shape.height",
      :x="shape.x + shape.width / 2",
      :editable="editing",
      align-center,
      justify-center,
      @blur="disableEditingMode"
    )
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
    readonly: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
