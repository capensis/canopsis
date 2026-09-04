<template>
  <c-layout gap="2" column>
    <c-enabled-field v-model="enabledField" :label="$t('common.enabled')" with-background />

    <c-form-general-patterns-tabs
      :form="values"
      :rule-id="ruleId"
      :type="type"
    >
      <template #general="{ setRef, templateVars, copyVars }">
        <dynamic-info-general-form
          :ref="setRef"
          :is-disabled-id-field="isDisabledIdField"
          :template-vars="templateVars"
          :copy-vars="copyVars"
        />
      </template>

      <template #patterns="{ setRef }">
        <dynamic-info-patterns-form :ref="setRef" />
      </template>
    </c-form-general-patterns-tabs>
  </c-layout>
</template>

<script>
import { TEMPLATE_TESTING_TEST_TYPES } from '@/constants';

import { useFormContext } from '@/hooks/form';

import DynamicInfoGeneralForm from './dynamic-info-general-form.vue';
import DynamicInfoPatternsForm from './dynamic-info-patterns-form.vue';

export default {
  components: {
    DynamicInfoGeneralForm,
    DynamicInfoPatternsForm,
  },

  props: {
    ruleId: {
      type: String,
      default: undefined,
    },
    isDisabledIdField: {
      type: Boolean,
      default: false,
    },
  },

  setup() {
    const form = useFormContext();
    const enabledField = form.model('enabled');
    const type = TEMPLATE_TESTING_TEST_TYPES.dynamicInfo;

    return {
      values: form.values,
      enabledField,
      type,
    };
  },
};
</script>
