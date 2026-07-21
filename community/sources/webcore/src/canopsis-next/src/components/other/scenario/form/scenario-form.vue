<template>
  <v-layout class="gap-2" column>
    <c-enabled-field v-field="form.enabled" with-background />
    <c-form-general-patterns-tabs
      v-field="form"
      :rule-id="ruleId"
      :type="type"
      :patterns-label="$t('common.actionsLabel')"
      :disabled-test-query-tooltip="disabledTestQueryTooltip"
    >
      <template #general="{ setRef }">
        <scenario-general-form v-field="form" :ref="setRef" />
      </template>

      <template #patterns="{ setRef, templateVars }">
        <scenario-actions-form v-field="form.actions" :ref="setRef" :template-vars="templateVars" />
      </template>

      <template #test-query>
        <scenario-test-query :form="form" />
      </template>
    </c-form-general-patterns-tabs>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { TEMPLATE_TESTING_TEST_TYPES } from '@/constants';

import { isWebhookActionType } from '@/helpers/entities/action/form';

import { useI18n } from '@/hooks/i18n';

import ScenarioTestQuery from '@/components/other/scenario/partials/scenario-test-query.vue';

import ScenarioGeneralForm from './scenario-general-form.vue';
import ScenarioActionsForm from './scenario-actions-form.vue';

export default {
  components: {
    ScenarioGeneralForm,
    ScenarioActionsForm,
    ScenarioTestQuery,
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
  },
  setup(props) {
    const { t } = useI18n();

    const type = TEMPLATE_TESTING_TEST_TYPES.scenario;

    const hasWebhookAction = computed(() => (
      props.form.actions.some(({ type: actionType }) => isWebhookActionType(actionType))
    ));

    const disabledTestQueryTooltip = computed(() => (hasWebhookAction.value ? '' : t('scenario.errors.testQueryRequireSteps')));

    return {
      type,
      disabledTestQueryTooltip,
    };
  },
};
</script>
