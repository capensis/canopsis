<template>
  <code-editor
    v-field="value"
    ref="codeEditor"
    :options="editorOptions"
    :resettable="resettable"
    language="javascript"
  />
</template>

<script>
import { computed, ref } from 'vue';

import { useJavaScriptCompletions } from '@/hooks/monaco';

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
    completions: {
      type: Object,
      required: false,
    },
  },
  setup(props) {
    const codeEditor = ref(null);

    const editorOptions = computed(() => ({
      theme: props.theme,
      automaticLayout: true,
      minimap: {
        enabled: false,
      },
    }));

    useJavaScriptCompletions({
      codeEditor,
      completions: props.completions,
    });

    return {
      codeEditor,

      editorOptions,
    };
  },
};
</script>
