<template lang="pug">
  div(ref="codeEditor")
</template>

<script>
import initEditor from 'monaco-mermaid';
import * as monaco from 'monaco-editor';

initEditor(monaco);

export default {
  mounted() {
    this.$editor = monaco.editor.create(this.$refs.codeEditor, {
      value: '',
      language: 'mermaid',
    });
  },
};
</script>
