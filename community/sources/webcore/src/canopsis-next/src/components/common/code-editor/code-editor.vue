<template>
  <div class="code-editor__wrapper">
    <c-action-btn
      v-if="resettable"
      :disabled="!wasChanged"
      :tooltip="$t('common.reset')"
      icon="$vuetify.icons.restart_alt"
      color="white"
      btn-color="grey darken-1"
      left
      @click="reset"
    />
    <div
      class="code-editor"
      ref="codeEditor"
    />
  </div>
</template>

<script>
import * as monaco from 'monaco-editor/esm/vs/editor/editor.api';
import initEditor from 'monaco-mermaid';
import iPlasticTheme from 'monaco-themes/themes/iPlastic.json';

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
    resettable: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    editorOptions() {
      return {
        ...this.options,

        language: this.language,
      };
    },

    wasChanged() {
      return this.originalValue !== this.value;
    },
  },
  watch: {
    value(value) {
      if (value !== this.$editor?.getValue()) {
        this.$editor.setValue(value);
        this.originalValue = value;
      }
    },

    editorOptions(options) {
      this.$editor?.updateOptions(options);
    },

    errorMarkers(errorMarkers) {
      if (this.$editor) {
        this.$monaco.editor.setModelMarkers(this.$editor.getModel(), 'test', errorMarkers);
      }
    },
  },
  created() {
    this.originalValue = this.value;
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
    this.$monaco.editor.defineTheme('iPlastic', iPlasticTheme);

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

    reset() {
      this.$emit('input', this.originalValue);
    },
  },
};
</script>

<style lang="scss" scoped>
.code-editor {
  position: absolute;
  height: 100%;
  width: 100%;

  &__wrapper {
    position: relative;

    ::v-deep .c-action-btn__button {
      position: absolute;
      right: 18px;
      top: 0;
      opacity: .6;
      z-index: 2;

      &:hover {
        opacity: 1;
      }
    }
  }
}
</style>
