<template>
  <v-layout class="gap-2" column>
    <c-enabled-field v-field="form.enabled" with-background />
    <c-form-general-patterns-tabs
      :form="form"
      :rule-id="ruleId"
      :type="type"
    >
      <template #general="{ setRef, templateVars, copyVars }">
        <dynamic-info-general-form
          v-field="form"
          :ref="setRef"
          :copy-vars="copyVars"
          :is-disabled-id-field="isDisabledIdField"
          :template-vars="templateVars"
        />
      </template>
      <template #patterns="{ setRef }">
        <dynamic-info-patterns-form v-field="form.patterns" :ref="setRef" />
      </template>
    </c-form-general-patterns-tabs>
  </v-layout>
</template>

<script>
import { TEMPLATE_TESTING_TEST_TYPES } from '@/constants';

import DynamicInfoGeneralForm from './dynamic-info-general-form.vue';
import DynamicInfoPatternsForm from './dynamic-info-patterns-form.vue';

export default {
  components: {
    DynamicInfoGeneralForm,
    DynamicInfoPatternsForm,
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
    isDisabledIdField: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    const type = TEMPLATE_TESTING_TEST_TYPES.dynamicInfo;

    return { type };
  },
};
</script>
