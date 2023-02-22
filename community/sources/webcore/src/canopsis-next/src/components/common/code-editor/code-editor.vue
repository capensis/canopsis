<template lang="pug">
  div(ref="codeEditor")
</template>

<script>
import * as monaco from 'monaco-editor/esm/vs/editor/editor.api';
import initEditor from 'monaco-mermaid';

import 'monaco-editor/esm/vs/basic-languages/javascript/javascript.contribution';

export default {
  props: {
    value: {
      type: String,
      default: '',
    },
    language: {
      type: String,
      default: 'mermaid',
    },
    options: {
      type: Object,
      default: () => ({}),
    },
    errorMarkers: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    editorOptions() {
      return {
        ...this.options,

        language: this.language,
      };
    },
  },
  watch: {
    editorOptions(value) {
      this.$editor?.updateOptions(value);
    },

    errorMarkers(errorMarkers) {
      if (this.$editor) {
        this.$monaco.editor.setModelMarkers(this.$editor.getModel(), 'test', errorMarkers);
      }
    },
  },
  async mounted() {
    try {
      this.$monaco = monaco;
    } catch {
      await this.loadMonaco();
    }

    if (!this.$monaco) {
      this.$popups.error({ text: this.$t('errors.codeEditorProblem') });

      return;
    }

    initEditor(this.$monaco);

    this.$editor = this.$monaco.editor.create(this.$refs.codeEditor, {
      value: this.value,

      ...this.editorOptions,
    });

    this.$editor.onDidChangeModelContent(this.onChange);
    this.$editor.onDidFocusEditorWidget(this.onFocus);
    this.$editor.onDidBlurEditorWidget(this.onBlur);
  },
  beforeDestroy() {
    this.$editor.dispose();
  },
  methods: {
    /**
     * This function fixes this problem: https://github.com/mermaid-js/mermaid-live-editor/issues/175
     *
     * @returns {Promise<void>}
     */
    async loadMonaco() {
      for (let i = 0; i < 10; i += 1) {
        try {
          this.$monaco = monaco;
          return;
        } catch {
          // eslint-disable-next-line no-await-in-loop
          await new Promise(resolve => setTimeout(resolve, 500));
        }
      }
    },

    onChange() {
      this.$emit('input', this.$editor.getValue());
    },

    onFocus(...args) {
      this.$emit('focus', ...args);
    },

    onBlur(...args) {
      this.$emit('blur', ...args);
    },
  },
};
</script>
