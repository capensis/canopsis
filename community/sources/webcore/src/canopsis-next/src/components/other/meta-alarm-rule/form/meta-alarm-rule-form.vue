<template>
  <v-layout class="gap-3" column>
    <c-enabled-field v-field="form.enabled" hide-details with-background />
    <c-form-general-patterns-tabs
      :form="form"
      :rule-id="ruleId"
      :type="type"
      :patterns-label="$t('metaAlarmRule.patternsTabLabel')"
    >
      <template #general="{ setRef, templateVars }">
        <meta-alarm-rule-general-form
          v-field="form"
          :ref="setRef"
          :disabled-id-field="disabledIdField"
          :template-vars="templateVars"
        />
      </template>
      <template #patterns="{ setRef, templateVars }">
        <meta-alarm-rule-parameters-form
          v-field="form"
          :ref="setRef"
          :template-vars="templateVars"
        />
      </template>
    </c-form-general-patterns-tabs>
  </v-layout>
</template>

<script>
import { TEMPLATE_TESTING_TEST_TYPES } from '@/constants';

import MetaAlarmRuleGeneralForm from './meta-alarm-rule-general-form.vue';
import MetaAlarmRuleParametersForm from './meta-alarm-rule-parameters-form.vue';

export default {
  components: {
    MetaAlarmRuleGeneralForm,
    MetaAlarmRuleParametersForm,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
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
    const type = TEMPLATE_TESTING_TEST_TYPES.metaAlarmRule;

    return { type };
  },
};
</script>
