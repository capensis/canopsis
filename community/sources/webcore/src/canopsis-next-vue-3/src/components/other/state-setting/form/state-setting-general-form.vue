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
        <c-select
          v-model="type"
          :label="t('stateSetting.appliedForEntityType')"
          :items="entityTypeItems"
          item-title="text"
          item-value="value"
          name="type"
          required
        />
      </c-form-block-row>

      <c-form-block-row :label="t('stateSetting.computeMethod')" indented>
        <state-setting-method-field v-model="method" />
      </c-form-block-row>
    </c-form-block>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { STATE_SETTING_ENTITY_TYPES } from '@/constants';

import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import StateSettingMethodField from './fields/state-setting-method-field.vue';

export default {
  components: { StateSettingMethodField },

  setup() {
    const { t, te } = useI18n();
    const { title, priority, type, method } = useFormModels([
      'title',
      'priority',
      'type',
      'method',
    ]);

    const entityTypeItems = computed(() => STATE_SETTING_ENTITY_TYPES.map(value => ({
        value,
        text: te(`common.${value}`) ? t(`common.${value}`) : value,
      })),
    );

    return {
      t,
      title,
      priority,
      type,
      method,
      entityTypeItems,
    };
  },
};
</script>
