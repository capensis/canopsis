<template lang="pug">
  code-editor(
    v-field="value",
    :options="editorOptions",
    language="javascript"
  )
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
  },
  data() {
    return {
      editorError: undefined,
    };
  },
  computed: {
    editorOptions() {
      return {
        automaticLayout: true,
        minimap: {
          enabled: false,
        },
      };
    },
  },
};
</script>
