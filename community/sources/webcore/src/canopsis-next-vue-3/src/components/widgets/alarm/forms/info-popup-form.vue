<template>
  <c-layout column>
    <c-select
      :model-value="form.column"
      :items="columns"
      :label="$t('common.column', 1)"
      name="column"
      item-title="text"
      item-value="value"
      autofocus
      required
      @update:model-value="updateField('column', $event)"
    />
    <text-editor-field
      :model-value="form.template"
      :variables="variables"
      name="template"
      required
      @update:model-value="updateField('template', $event)"
    />
  </c-layout>
</template>

<script>
import { useModelField } from '@/hooks/form/model-field';

import TextEditorField from '@/components/forms/fields/text-editor-field.vue';

export default {
  components: { TextEditorField },

  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    columns: {
      type: Array,
      default: () => [],
    },
    variables: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['update:form'],

  setup(props, { emit }) {
    const { updateField } = useModelField(props, emit);

    return {
      updateField,
    };
  },
};
</script>
