<template lang="pug">
  v-layout(column)
    c-name-field(v-field="form.name", required)
    mermaid-editor(v-field="form.parameters")
</template>

<script>
import MermaidEditor from './fields/mermaid-editor.vue';

export default {
  inject: ['$validator'],
  components: { MermaidEditor },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
  },
};
</script>
