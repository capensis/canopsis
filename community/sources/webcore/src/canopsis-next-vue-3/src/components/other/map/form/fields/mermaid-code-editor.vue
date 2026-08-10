<template>
  <code-editor
    :value="value"
    :options="editorOptions"
    :error-markers="errorMarkers"
    @update:value="updateModel"
  />
</template>

<script>
import CodeEditor from '@/components/common/code-editor/code-editor.vue';
import { validateMermaidDiagram } from '@/helpers/mermaid';
import { formBaseMixin } from '@/mixins/form';

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
