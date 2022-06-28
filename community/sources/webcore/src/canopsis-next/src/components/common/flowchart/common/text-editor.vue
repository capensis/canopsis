<template lang="pug">
  component(is="foreignObject", width="100%", height="100%", pointer-events="none")
    div.text-shape-editor(:class="editorClasses", :style="editorStyle")
      div.text-shape-editor__wrapper
        div.text-shape-editor__field(
          ref="textEditor",
          v-html="value",
          :contenteditable="editable",
          :style="fieldStyle",
          :class="fieldClass",
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
      required: false,
    },
    height: {
      type: Number,
      required: false,
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
    hidden: {
      type: Boolean,
      default: false,
    },
    fieldClass: {
      type: String,
      required: false,
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

    editorStyle() {
      return {
        top: `${this.y}px`,
        left: `${this.x}px`,
        width: this.width ? `${this.width}px` : '1px',
        height: this.height ? `${this.height}px` : '1em',
        overflow: !this.editable && this.hidden ? 'hidden' : '',
      };
    },

    fieldStyle() {
      return {
        width: this.width ? `${this.width}px` : 'max-content',
        'text-align': this.justifyCenter ? 'center' : '',
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
  position: absolute;
  display: flex;
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
