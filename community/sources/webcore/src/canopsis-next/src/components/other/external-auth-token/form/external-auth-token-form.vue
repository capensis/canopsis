<template>
  <c-form-general-patterns-tabs
    :form="formValues"
    :rule-id="ruleId"
    :type="type"
    @update:form="updateForm"
  >
    <template #general="{ setRef, templateVars }">
      <external-auth-token-general-form :ref="setRef" :template-vars="templateVars" />
    </template>
  </c-form-general-patterns-tabs>
</template>

<script>
import { TEMPLATE_TESTING_TEST_TYPES } from '@/constants';

import { useFormContext } from '@/hooks/form';

import ExternalAuthTokenGeneralForm from './external-auth-token-general-form.vue';

export default {
  components: { ExternalAuthTokenGeneralForm },

  props: {
    ruleId: {
      type: String,
      required: false,
    },
  },

  setup() {
    const form = useFormContext();
    const type = TEMPLATE_TESTING_TEST_TYPES.externalAuthToken;

    const updateForm = (value) => {
      form.values.value = value;
    };

    return {
      formValues: form.values,
      type,
      updateForm,
    };
  },
};
</script>
