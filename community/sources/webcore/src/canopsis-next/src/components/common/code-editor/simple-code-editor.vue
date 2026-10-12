<template>
  <code-editor
    v-field="value"
    ref="codeEditor"
    :options="editorOptions"
    :resettable="resettable"
    :language="language"
  />
</template>

<script>
import { computed, ref } from 'vue';

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
    theme: {
      type: String,
      default: 'iPlastic',
    },
    resettable: {
      type: Boolean,
      default: false,
    },
    readonly: {
      type: Boolean,
      default: false,
    },
    lineNumbers: {
      type: String,
      default: 'off',
    },
    renderLineHighlight: {
      type: String,
      default: 'none',
    },
    minimap: {
      type: Boolean,
      default: false,
    },
    language: {
      type: String,
      default: 'plaintext',
    },
  },
  setup(props) {
    const codeEditor = ref(null);

    const editorOptions = computed(() => ({
      theme: props.theme,
      automaticLayout: true,
      lineNumbers: props.lineNumbers,
      readOnly: props.readonly,
      foldingStrategy: 'indentation',
      renderLineHighlight: props.renderLineHighlight,
      minimap: {
        enabled: props.minimap,
      },
    }));

    return {
      codeEditor,

      editorOptions,
    };
  },
};
</script>
