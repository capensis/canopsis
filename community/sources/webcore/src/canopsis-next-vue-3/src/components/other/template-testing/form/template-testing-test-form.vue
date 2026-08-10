<template>
  <c-layout gap="2" column>
    <c-name-field
      v-model="nameField"
      :label="t('common.testName')"
      name="name"
      required
    />

    <template-testing-test-type-field
      :value="typeField"
      :label="t('common.ruleType')"
      :disabled="!isNew"
      name="type"
      required
      @update:value="typeField = $event"
    />

    <c-name-field
      :model-value="ruleField?.name"
      :label="t('common.ruleName')"
      :disabled="!isNew"
      name="rule"
      required
    />

    <c-description-field
      v-model="descriptionField"
      :label="t('common.description')"
      :max-length="500"
      name="description"
    />
  </c-layout>
</template>

<script>
import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import TemplateTestingTestTypeField from './fields/template-testing-test-type-field.vue';

export default {
  components: {
    TemplateTestingTestTypeField,
  },
  props: {
    isNew: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    const { t } = useI18n();
    const { nameField, typeField, ruleField, descriptionField } = useFormModels({
      nameField: 'name',
      typeField: 'type',
      ruleField: 'rule',
      descriptionField: 'description',
    });

    return {
      t,
      nameField,
      typeField,
      ruleField,
      descriptionField,
    };
  },
};
</script>
