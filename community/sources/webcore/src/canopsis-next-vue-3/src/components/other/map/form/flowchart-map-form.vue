<template>
  <c-layout class="geomap-map-form" column>
    <c-name-field
      :model-value="form.name"
      autofocus
      required
      @update:model-value="updateField('name', $event)"
    />
    <flowchart-editor :form="form.parameters" @update:form="updateField('parameters', $event)" />
  </c-layout>
</template>

<script>
import { formMixin } from '@/mixins/form';

import FlowchartEditor from './fields/flowchart-editor.vue';

export default {
  components: { FlowchartEditor },
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

<style lang="scss">
.geomap-map-form {
  &__editor {
    height: 500px;
  }
}
</style>
