<template>
  <div class="position-relative">
    <c-progress-overlay :pending="pending" />
    <v-tabs fixed-tabs>
      <v-tab>{{ $t('common.general') }}</v-tab>
      <template-testing-test-variables-tab :disabled="isEmptyVariablesFields" />

      <v-tab-item class="pt-2" eager>
        <slot :template-vars="templateVars" />
      </v-tab-item>

      <v-tab-item :disabled="isEmptyVariablesFields">
        <template-testing-test-variables
          :general-form="form"
          :variables-fields="variablesFields"
          :template-vars="templateVars"
          :rule-id="ruleId"
          :type="type"
        />
      </v-tab-item>
    </v-tabs>
  </div>
</template>

<script>
import { onMounted, toRef } from 'vue';

import { useTemplateVarsList, useTestVariablesFields } from './hooks/template-test-variables-wrapper';
import TemplateTestingTestVariables from './template-testing-test-variables.vue';
import TemplateTestingTestVariablesTab from './template-testing-test-variables-tab.vue';

export default {
  components: { TemplateTestingTestVariables, TemplateTestingTestVariablesTab },
  inheritAttrs: false,
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
    type: {
      type: Number,
      required: false,
    },
  },
  setup(props, { emit }) {
    const { templateVars, pending, fetchList } = useTemplateVarsList({
      type: toRef(props, 'type'),
      form: toRef(props, 'form'),
    });

    const {
      items: variablesFields,
      isEmptyItems: isEmptyVariablesFields,
    } = useTestVariablesFields(props, toRef(props, 'type'), emit);

    onMounted(fetchList);

    return {
      variablesFields,
      templateVars,
      pending,
      isEmptyVariablesFields,
    };
  },
};
</script>
