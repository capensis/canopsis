<template>
  <g
    class="image-shape"
    @dblclick.stop="enableEditingMode"
  >
    <template v-if="shape.svg">
      <component
        is="foreignObject"
        v-on="$listeners"
        v-html="shape.svg"
        :style="svgStyle"
        :x="shape.x"
        :y="shape.y"
        :height="shape.height"
        :width="shape.width"
        cursor="move"
      />
    </template>
    <template v-else>
      <rect
        v-bind="shape.properties"
        :x="shape.x"
        :y="shape.y"
        :width="shape.width"
        :height="shape.height"
      />
      <image
        v-on="$listeners"
        :href="shape.src"
        :x="shape.x"
        :y="shape.y"
        :width="shape.width"
        :height="shape.height"
        :cursor="readonly ? '' : 'move'"
      />
    </template>
    <text-editor
      ref="editor"
      v-bind="shape.textProperties"
      :value="shape.text"
      :y="shape.y + shape.height + textOffsetY"
      :x="shape.x + shape.width / 2"
      :editable="editing"
      align-center
      justify-center
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
    readonly: {
      type: Boolean,
      default: false,
    },
    textOffsetY: {
      type: Number,
      default: 5,
    },
  },
  computed: {
    svgStyle() {
      return {
        color: this.shape.properties.fill,
      };
    },
  },
};
</script>

<style lang="scss">
.image-shape {
  svg {
    width: 100%;
    height: 100%;
  }
}
</style>
