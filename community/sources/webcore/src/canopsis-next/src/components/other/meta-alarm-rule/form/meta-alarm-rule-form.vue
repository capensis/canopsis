<template>
  <c-layout gap="3" column>
    <c-enabled-field v-model="enabledField" hide-details with-background />

    <c-form-general-patterns-tabs
      :form="formValues"
      :rule-id="ruleId"
      :type="type"
      :patterns-label="t('metaAlarmRule.patternsTabLabel')"
    >
      <template #general="{ setRef, templateVars }">
        <meta-alarm-rule-general-form
          :ref="setRef"
          :disabled-id-field="disabledIdField"
          :template-vars="templateVars"
        />
      </template>

      <template #patterns="{ setRef, templateVars }">
        <meta-alarm-rule-parameters-form
          :ref="setRef"
          :template-vars="templateVars"
        />
      </template>
    </c-form-general-patterns-tabs>
  </c-layout>
</template>

<script>
import { TEMPLATE_TESTING_TEST_TYPES } from '@/constants';

import { useFormContext } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import MetaAlarmRuleGeneralForm from './meta-alarm-rule-general-form.vue';
import MetaAlarmRuleParametersForm from './meta-alarm-rule-parameters-form.vue';

export default {
  components: {
    MetaAlarmRuleGeneralForm,
    MetaAlarmRuleParametersForm,
  },

  props: {
    ruleId: {
      type: String,
      default: undefined,
    },
    disabledIdField: {
      type: Boolean,
      default: false,
    },
  },

  setup() {
    const { t } = useI18n();
    const form = useFormContext();
    const enabledField = form.model('enabled');
    const type = TEMPLATE_TESTING_TEST_TYPES.metaAlarmRule;

    return {
      t,
      formValues: form.values,
      enabledField,
      type,
    };
  },
};
</script>
