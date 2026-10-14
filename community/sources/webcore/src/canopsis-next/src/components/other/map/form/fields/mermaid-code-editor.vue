<template>
  <code-editor
    :value="value"
    :options="editorOptions"
    :error-markers="errorMarkers"
    @update:value="updateModel"
  />
</template>

<script>
import { validateMermaidDiagram } from '@/helpers/mermaid';

import { formBaseMixin } from '@/mixins/form';

import CodeEditor from '@/components/common/code-editor/code-editor.vue';

export default {
  components: { CodeEditor },
  mixins: [formBaseMixin],
  props: {
    value: {
      type: String,
      default: '',
    },
  },
  emits: ['input'],
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
    async value(value) {
      try {
        await validateMermaidDiagram(value);

        this.editorError = undefined;
      } catch (error) {
        const location = error.hash?.loc ?? {
          first_line: 1,
          first_column: 0,
          last_line: 1,
          last_column: 0,
        };

        this.editorError = {
          severity: 8,
          startLineNumber: location.first_line,
          startColumn: location.first_column,
          endLineNumber: location.last_line,
          endColumn: location.last_column + 1,
          message: error.str ?? error.message,
        };
      }
    },
  },
};
</script>
