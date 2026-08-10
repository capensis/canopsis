<template>
  <c-layout column>
    <c-select
      :model-value="form.column"
      :items="preparedColumns"
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
import { computed } from 'vue';

import TextEditorField from '@/components/forms/fields/text-editor-field.vue';
import { ALARM_FIELDS_TO_LABELS_KEYS } from '@/constants';
import { formToWidgetColumns } from '@/helpers/entities/widget/column/form';
import { getWidgetColumnLabel } from '@/helpers/entities/widget/list';
import { useModelField } from '@/hooks/form/model-field';

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
    const preparedColumns = computed(() =>
      formToWidgetColumns(props.columns).map(column => ({
        text: getWidgetColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
        value: column.value,
      })),
    );

    return {
      preparedColumns,
      updateField,
    };
  },
};
</script>
