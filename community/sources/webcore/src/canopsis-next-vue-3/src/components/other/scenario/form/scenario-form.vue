<template>
  <c-layout gap="2" column>
    <c-enabled-field v-model="enabledField" with-background />

    <c-form-general-patterns-tabs
      :form="values"
      :rule-id="ruleId"
      :type="testType"
      :patterns-label="t('common.actionsLabel')"
      :disabled-test-query-tooltip="disabledTestQueryTooltip"
    >
      <template #general="{ setRef }">
        <scenario-general-form :ref="setRef" />
      </template>

      <template #patterns="{ setRef, templateVars }">
        <scenario-actions-form :ref="setRef" :template-vars="templateVars" />
      </template>

      <template #test-query="{ setRef }">
        <scenario-test-query :ref="setRef" :form="values" />
      </template>
    </c-form-general-patterns-tabs>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { TEMPLATE_TESTING_TEST_TYPES } from '@/constants';

import { isWebhookActionType } from '@/helpers/entities/action/form';

import { useFormContext } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import ScenarioTestQuery from '@/components/other/scenario/partials/scenario-test-query.vue';

import ScenarioActionsForm from './scenario-actions-form.vue';
import ScenarioGeneralForm from './scenario-general-form.vue';

export default {
  components: {
    ScenarioActionsForm,
    ScenarioGeneralForm,
    ScenarioTestQuery,
  },

  props: {
    ruleId: {
      type: String,
      required: false,
    },
  },

  setup() {
    const { t } = useI18n();
    const form = useFormContext();
    const enabledField = form.model('enabled');
    const hasWebhookAction = computed(() => form.values.value.actions?.some(action => isWebhookActionType(action.type)),
    );
    const disabledTestQueryTooltip = computed(() => hasWebhookAction.value ? '' : t('scenario.errors.testQueryRequireSteps'),
    );

    return {
      t,
      values: form.values,
      enabledField,
      testType: TEMPLATE_TESTING_TEST_TYPES.scenario,
      disabledTestQueryTooltip,
    };
  },
};
</script>
