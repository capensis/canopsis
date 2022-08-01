<template lang="pug">
  v-layout(column)
    v-text-field(
      v-field="form.name",
      v-validate="'required'",
      :label="$t('common.name')",
      :error-messages="errors.collect('name')",
      name="name"
    )
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
