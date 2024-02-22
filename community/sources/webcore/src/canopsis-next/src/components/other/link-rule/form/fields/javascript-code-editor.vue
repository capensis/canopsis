<template>
  <code-editor
    v-field="value"
    :options="editorOptions"
    :resettable="resettable"
    language="javascript"
  />
</template>

<script>
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
  data() {
    return {
      editorError: undefined,
    };
  },
  computed: {
    editorOptions() {
      return {
        theme: this.theme,
        automaticLayout: true,
        minimap: {
          enabled: false,
        },
      };
    },
  },
};
</script>
