<template>
  <v-tabs
    slider-color="primary"
    fixed-tabs
  >
    <v-tab :class="{ 'error--text': hasGeneralError }">
      {{ $t('common.general') }}
    </v-tab>
    <v-tab :class="{ 'error--text': hasPatternsError }">
      {{ $tc('common.pattern', 2) }}
    </v-tab>

    <template-testing-test-variables-tab :disabled="isEmptyVariablesFields" />

    <v-tab-item eager>
      <remediation-instruction-general-form
        v-field="form"
        ref="generalElement"
        :disabled="disabled"
        :is-new="isNew"
        :required-approve="requiredApprove"
        :template-vars="templateVars"
        class="mt-3"
      />
    </v-tab-item>
    <v-tab-item eager>
      <remediation-instruction-patterns-form
        v-field="form.patterns"
        ref="patternsElement"
        class="mt-3"
      />
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
</template>

<script>
import {
  ref,
  toRef,
  watch,
  onMounted,
  onBeforeUnmount,
} from 'vue';

import { TEMPLATE_TESTING_TEST_TYPES } from '@/constants';

import { useTemplateVarsList } from '@/hooks/vars/template';

import { useTestVariablesFields } from '@/components/other/template-testing/hooks/template-test-variables-wrapper';

import TemplateTestingTestVariables from '@/components/other/template-testing/template-testing-test-variables.vue';
import TemplateTestingTestVariablesTab from '@/components/other/template-testing/template-testing-test-variables-tab.vue';

import RemediationInstructionGeneralForm from './remediation-instruction-general-form.vue';
import RemediationInstructionPatternsForm from './remediation-instruction-patterns-form.vue';

export default {
  inject: ['$validator'],
  components: {
    TemplateTestingTestVariables,
    TemplateTestingTestVariablesTab,

    RemediationInstructionGeneralForm,
    RemediationInstructionPatternsForm,
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
    disabled: {
      type: Boolean,
      default: false,
    },
    disabledCommon: {
      type: Boolean,
      default: false,
    },
    isNew: {
      type: Boolean,
      default: false,
    },
    requiredApprove: {
      type: Boolean,
      default: false,
    },
    ruleId: {
      type: String,
      required: false,
    },
  },
  setup(props, { emit }) {
    const hasGeneralError = ref(false);
    const hasPatternsError = ref(false);

    const generalElement = ref(null);
    const patternsElement = ref(null);

    const type = ref(TEMPLATE_TESTING_TEST_TYPES.instruction);

    const { vars: templateVars, fetchList } = useTemplateVarsList({
      type,
      form: toRef(props, 'form'),
    });

    const {
      items: variablesFields,
      isEmptyItems: isEmptyVariablesFields,
    } = useTestVariablesFields(props, type, emit);

    let unwatchGeneralTabErrors = null;
    let unwatchPatternsTabErrors = null;

    const watchTabsErrors = () => {
      unwatchGeneralTabErrors = watch(() => generalElement.value?.hasAnyError, (value) => {
        hasGeneralError.value = value;
      });

      unwatchPatternsTabErrors = watch(() => patternsElement.value?.hasAnyError, (value) => {
        hasPatternsError.value = value;
      });
    };

    const unwatchTabsErrors = () => {
      unwatchGeneralTabErrors?.();
      unwatchPatternsTabErrors?.();
    };

    onMounted(() => {
      watchTabsErrors();
      fetchList();
    });

    onBeforeUnmount(unwatchTabsErrors);

    return {
      generalElement,
      patternsElement,
      hasGeneralError,
      hasPatternsError,
      variablesFields,
      templateVars,
      isEmptyVariablesFields,
      type,
    };
  },
};
</script>
