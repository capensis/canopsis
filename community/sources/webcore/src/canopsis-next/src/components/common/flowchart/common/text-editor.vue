<template>
  <component
    :is="'foreignObject'"
    pointer-events="none"
    style="overflow: visible;"
  >
    <div
      class="text-shape-editor"
      :class="editorClasses"
      :style="editorStyle"
    >
      <div class="text-shape-editor__wrapper">
        <div
          class="text-shape-editor__field"
          ref="textEditor"
          v-html="value"
          :contenteditable="editable"
          :style="fieldStyle"
          @mousedown.stop=""
          @mouseup.stop=""
          @keydown.stop=""
          @blur="$emit('blur', $event)"
        />
      </div>
    </div>
  </component>
</template>

<script>
export default {
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
    color: {
      type: String,
      default: 'black',
    },
    fontSize: {
      type: Number,
      default: 12,
    },
    backgroundColor: {
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
        color: this.color,
        backgroundColor: this.backgroundColor,
        fontSize: `${this.fontSize}px`,
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

  &--editable {
    pointer-events: all;
  }

  &:not(&--editable) {
    pointer-events: none;
  }
}
</style>
