<template>
  <c-form-general-patterns-tabs
    v-field="form"
    :rule-id="ruleId"
    :type="type"
  >
    <template #general="{ setRef, templateVars }">
      <remediation-job-general-form
        v-field="form"
        :ref="setRef"
        :template-vars="templateVars"
      />
    </template>
  </c-form-general-patterns-tabs>
</template>

<script>
import { TEMPLATE_TESTING_TEST_TYPES } from '@/constants';

import RemediationJobGeneralForm from './remediation-job-general-form.vue';

export default {
  inject: ['$validator'],
  components: {
    RemediationJobGeneralForm,
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
      required: false,
    },
  },
  setup() {
    const type = TEMPLATE_TESTING_TEST_TYPES.job;

    return {
      type,
    };
  },
};
</script>
