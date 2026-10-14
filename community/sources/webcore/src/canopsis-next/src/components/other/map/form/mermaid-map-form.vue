<template>
  <c-layout column>
    <c-name-field
      :model-value="form.name"
      autofocus
      required
      @update:model-value="updateField('name', $event)"
    />
    <mermaid-editor :form="form.parameters" @update:form="updateField('parameters', $event)" />
  </c-layout>
</template>

<script>
import { formMixin } from '@/mixins/form';

import MermaidEditor from './fields/mermaid-editor.vue';

export default {
  components: { MermaidEditor },
  mixins: [formMixin],
  props: {
    form: {
      type: Object,
      required: true,
    },
  },
  emits: ['update:form'],
};
</script>
