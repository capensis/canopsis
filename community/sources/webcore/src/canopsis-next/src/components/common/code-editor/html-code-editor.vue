<template>
  <code-editor
    v-field="value"
    ref="codeEditor"
    :options="editorOptions"
    :resettable="resettable"
    language="html"
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

    return {
      codeEditor,

      editorOptions,
    };
  },
};
</script>
