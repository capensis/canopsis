<template lang="pug">
  component(is="foreignObject", :y="y", :x="x", :width="width", :height="height")
    div.text-shape-editor(:class="editorClasses")
      div.text-shape-editor__field(
        ref="textEditor",
        v-html="value",
        :contenteditable="editable",
        @mousedown.stop="",
        @mouseup.stop="",
        @keydown.stop="",
        @blur="$emit('blur', $event)"
      )
</template>

<script>
import { formBaseMixin } from '@/mixins/form';

import RectShapeSelection from '../rect-shape/rect-shape-selection.vue';

export default {
  components: { RectShapeSelection },
  mixins: [formBaseMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: '',
    },
    x: {
      type: Number,
      required: true,
    },
    y: {
      type: Number,
      required: true,
    },
    width: {
      type: Number,
      required: true,
    },
    height: {
      type: Number,
      required: true,
    },
    editable: {
      type: Boolean,
      default: false,
    },
    alignCenter: {
      type: Boolean,
      default: false,
    },
    justifyCenter: {
      type: Boolean,
      default: false,
    },
    unselect: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    editorClasses() {
      return {
        'text-shape-editor--align-center': this.alignCenter,
        'text-shape-editor--justify-center': this.justifyCenter,
        'text-shape-editor--unselect': this.unselect,
        'text-shape-editor--editable': this.editable,
      };
    },
  },
  methods: {
    focus() {
      const { textEditor } = this.$refs;

      if (this.value) {
        const range = document.createRange();
        const selection = window.getSelection();

        range.selectNodeContents(textEditor);

        selection.removeAllRanges();
        selection.addRange(range);
      }

      this.$nextTick(() => textEditor.focus());
    },
  },
};
</script>

<style lang="scss" scoped>
.text-shape-editor {
  display: flex;
  height: 100%;
  width: 100%;
  user-select: none;

  &__field {
    min-width: 1px;
    min-height: 1em;

    &:focus {
      outline: none;
    }
  }

  &--align-center {
    align-items: center;
  }

  &--justify-center {
    justify-content: center;
  }

  &:not(&--editable) {
    pointer-events: none;
  }
}
</style>
