<template>
  <c-layout class="mt-3" gap="2" column>
    <span class="text-body-2">
      {{ t('common.column', 2) }}
    </span>

    <c-columns-field
      v-model="columns"
      :type="entityType"
      name="columns"
      with-color-indicator
      with-template
      with-html
    />
  </c-layout>
</template>

<script>
import { computed, watch } from 'vue';

import { ENTITIES_TYPES, WIDGET_TEMPLATES_TYPES } from '@/constants';

import { useField, useFormContext } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

export default {
  setup() {
    const { t } = useI18n();
    const form = useFormContext();
    const columns = form.model('columns');
    const type = form.model('type');

    const { validate } = useField(
      'columns',
      () => (columns.value.length > 0 ? true : t('widgetTemplate.errors.columnsRequired')),
      {
        initialValue: columns.value.length,
        validateOnValueUpdate: false,
      },
    );

    const entityType = computed(() =>
      type.value === WIDGET_TEMPLATES_TYPES.alarmColumns
        ? ENTITIES_TYPES.alarm
        : ENTITIES_TYPES.entity,
    );

    watch(
      () => columns.value.length,
      () => validate(),
      { flush: 'post' },
    );

    return {
      t,
      columns,
      entityType,
    };
  },
};
</script>
