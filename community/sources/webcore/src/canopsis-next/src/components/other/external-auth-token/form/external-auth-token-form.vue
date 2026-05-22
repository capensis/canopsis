<template>
  <c-form-general-patterns-tabs
    :form="form"
    :rule-id="ruleId"
    :type="type"
  >
    <template #general="{ setRef, templateVars }">
      <external-auth-token-general-form
        v-field="form"
        :ref="setRef"
        :template-vars="templateVars"
      />
    </template>
  </c-form-general-patterns-tabs>
</template>

<script>
import { TEMPLATE_TESTING_TEST_TYPES } from '@/constants';

import ExternalAuthTokenGeneralForm from './external-auth-token-general-form.vue';

export default {
  components: { ExternalAuthTokenGeneralForm },
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
    const type = TEMPLATE_TESTING_TEST_TYPES.externalAuthToken;

    return {
      type,
    };
  },
};
</script>
