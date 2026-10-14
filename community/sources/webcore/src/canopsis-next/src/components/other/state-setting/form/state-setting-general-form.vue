<template>
  <c-layout gap="3" column>
    <c-name-field
      v-model="title"
      :label="t('common.title')"
      name="title"
      autofocus
      required
    />

    <c-form-block>
      <c-form-block-row :label="t('common.priority')">
        <c-priority-field v-model="priority" required />
      </c-form-block-row>

      <c-form-block-row :label="t('stateSetting.appliedForEntityType')">
        <state-setting-type-field v-model="type" />
      </c-form-block-row>

      <c-form-block-row :label="t('stateSetting.computeMethod')">
        <state-setting-method-field v-model="method" />
      </c-form-block-row>
    </c-form-block>
  </c-layout>
</template>

<script>
import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import StateSettingMethodField from './fields/state-setting-method-field.vue';
import StateSettingTypeField from './fields/state-setting-type-field.vue';

export default {
  components: { StateSettingMethodField, StateSettingTypeField },

  setup() {
    const { t } = useI18n();
    const { title, priority, type, method } = useFormModels([
      'title',
      'priority',
      'type',
      'method',
    ]);

    return {
      t,
      title,
      priority,
      type,
      method,
    };
  },
};
</script>
