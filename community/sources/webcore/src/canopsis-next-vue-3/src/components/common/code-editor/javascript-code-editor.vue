<template>
  <code-editor
    ref="codeEditor"
    :value="value"
    :options="editorOptions"
    :resettable="resettable"
    language="javascript"
    @update:value="$emit('update:value', $event)"
  />
</template>

<script>
import { computed, ref, toRef } from 'vue';

import CodeEditor from '@/components/common/code-editor/code-editor.vue';
import { useJavaScriptCompletions } from '@/hooks/monaco';

export default {
  components: { CodeEditor },
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

  emits: ['update:value'],

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
      completions: toRef(props, 'completions'),
    });

    return {
      codeEditor,

      editorOptions,
    };
  },
};
</script>
