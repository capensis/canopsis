<template>
  <code-editor
    v-field="value"
    :options="editorOptions"
    :error-markers="errorMarkers"
  />
</template>

<script>
import { validateMermaidDiagram } from '@/helpers/mermaid';

import CodeEditor from '@/components/common/code-editor/code-editor.vue';

export default {
  components: { CodeEditor },
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: '',
    },
  },
  data() {
    return {
      editorError: undefined,
    };
  },
  computed: {
    editorOptions() {
      return {
        overviewRulerLanes: 0,
        theme: 'mermaid',
        automaticLayout: true,
        minimap: {
          enabled: false,
        },
      };
    },

    errorMarkers() {
      return this.editorError ? [this.editorError] : [];
    },
  },
  watch: {
    value(value) {
      try {
        validateMermaidDiagram(value);

        this.editorError = undefined;
      } catch (error) {
        this.editorError = {
          severity: 8,
          startLineNumber: error.hash.loc.first_line,
          startColumn: error.hash.loc.first_column,
          endLineNumber: error.hash.loc.last_line,
          endColumn: error.hash.loc.last_column + 1,
          message: error.str,
        };
      }
    },
  },
};
</script>
