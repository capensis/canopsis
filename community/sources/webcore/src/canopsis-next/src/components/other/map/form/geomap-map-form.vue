<template>
  <c-layout column>
    <c-name-field
      :model-value="form.name"
      autofocus
      required
      @update:model-value="updateField('name', $event)"
    />
    <geomap-editor :form="form.parameters" @update:form="updateField('parameters', $event)" />
  </c-layout>
</template>

<script>
import { formMixin } from '@/mixins/form';

import GeomapEditor from './fields/geomap-editor.vue';

export default {
  components: { GeomapEditor },
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
