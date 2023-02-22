<template lang="pug">
  java-script-code-editor.fill-min-height(v-field="value")
</template>

<script>
import JavaScriptCodeEditor from './partials/javascript-code-editor.vue';

export default {
  components: { JavaScriptCodeEditor },
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
};
</script>
