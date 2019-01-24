<template lang="pug">
  div.text-editor-wrapper
    quill-editor(
    :value="value",
    :disabled="disabled",
    @input="$emit('input', $event)"
    )
    codemirror.html-editor(
    v-if="!disabled && buttonToggle === 0",
    v-model="code",
    :options="codeMirrorOptions",
    @blur="$emit('input', code)"
    )
    v-btn-toggle.toggle-button(v-show="!disabled", v-model="buttonToggle")
      v-btn(flat, small)
        v-icon(small) code

</template>

<script>
import { quillEditor as QuillEditor } from 'vue-quill-editor';
import beautify from 'js-beautify';

import 'codemirror/mode/xml/xml';
import 'codemirror/addon/selection/active-line';
import 'codemirror/addon/edit/closetag';

import 'codemirror/theme/idea.css';


export default {
  components: {
    QuillEditor,
  },
  props: {
    value: {
      type: String,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      code: this.value,
      buttonToggle: null,
      codeMirrorOptions: {
        tabSize: 2,
        styleActiveLine: true,
        lineNumbers: false,
        autoCloseTags: true,
        line: true,
        mode: 'text/html',
        theme: 'idea',
      },
    };
  },
  watch: {
    buttonToggle(value) {
      if (value === 0) {
        /**
         * Code tab is active
         */
        this.parseHtmlToCode();
      } else {
        /**
         * Code tab is inactive
         */
        this.$emit('input', this.code);
      }
    },
  },
  methods: {
    parseHtmlToCode() {
      const beautifyOptions = {
        indent_size: 2,
        html: {
          end_with_newline: true,
          js: { indent_size: 2 },
          css: { indent_size: 2 },
        },
        css: { indent_size: 1 },
      };

      /**
       * Reformat HTML code from quill
       */
      this.code = beautify.html(this.value, beautifyOptions);
    },
  },
};
</script>

<style lang="scss">
  .text-editor-wrapper {
    position: relative;

    .html-editor {
      cursor: text;
      border: 1px solid #ccc;
      position: absolute;
      top: 0;
      left: 0;
      height: 100%;
      width: 100%;

      .vue-codemirror, .CodeMirror {
        height: 100%;
      }
    }

    .toggle-button {
      position: absolute;
      right: 1px;
      bottom: 1px;

      &.v-btn-toggle--selected {
        box-shadow: none;
      }
    }
  }
</style>

