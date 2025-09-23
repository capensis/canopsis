<template>
  <div class="position-relative">
    <c-progress-overlay :pending="pending" />
    <v-tabs fixed-tabs>
      <v-tab>{{ $t('common.general') }}</v-tab>
      <template-testing-test-variables-tab :disabled="isEmptyVariablesFields" />

      <v-tab-item class="pt-2" eager>
        <slot :template-vars="templateVars" :copy-vars="copyVars" />
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
import { computed, onMounted, toRef } from 'vue';

import { useCopyVarsList } from '@/hooks/vars/copy';
import { useTemplateVarsList } from '@/hooks/vars/template';

import { useTestVariablesFields } from './hooks/template-test-variables-wrapper';
import TemplateTestingTestVariables from './template-testing-test-variables.vue';
import TemplateTestingTestVariablesTab from './partials/template-testing-test-variables-tab.vue';

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
    const {
      vars: copyVars,
      pending: copyVarsPending,
      fetchList: fetchCopyVarsList,
    } = useCopyVarsList({
      type: toRef(props, 'type'),
    });

    const {
      vars: templateVars,
      pending: templateVarsPending,
      fetchList: fetchTemplateVarsList,
    } = useTemplateVarsList({
      type: toRef(props, 'type'),
      form: toRef(props, 'form'),
    });

    const {
      items: variablesFields,
      isEmptyItems: isEmptyVariablesFields,
    } = useTestVariablesFields(props, toRef(props, 'type'), emit);

    const pending = computed(() => templateVarsPending.value || copyVarsPending.value);

    onMounted(() => {
      fetchCopyVarsList();
      fetchTemplateVarsList();
    });

    return {
      copyVars,
      templateVars,
      pending,

      variablesFields,
      isEmptyVariablesFields,
    };
  },
};
</script>
