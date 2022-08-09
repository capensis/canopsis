<template lang="pug">
  v-layout(column)
    c-name-field(v-field="form.name")
    mermaid-editor(v-model="form.properties")
</template>

<script>
import MermaidEditor from '@/components/common/mermaid/mermaid-editor.vue';

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
