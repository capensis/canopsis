<template>
  <v-layout class="gap-3" column>
    <c-enabled-field
      v-field="form.enabled"
      :disabled="disabledCommon"
      with-background
    />

    <c-form-general-patterns-tabs
      :form="form"
      :rule-id="ruleId"
      :type="type"
      :additional-label="additionalLabel"
    >
      <template #general="{ setRef, templateVars }">
        <remediation-instruction-general-form
          v-field="form"
          :ref="setRef"
          :disabled="disabled"
          :disabled-common="disabledCommon"
          :is-new="isNew"
          :required-approve="requiredApprove"
          :template-vars="templateVars"
        />
      </template>

      <template #additional="{ setRef, templateVars }">
        <remediation-instruction-steps-form
          v-if="isManualType"
          v-field="form.steps"
          :ref="setRef"
          :disabled="disabled"
          :template-vars="templateVars"
          class="mt-3"
        />
        <remediation-instruction-jobs-form
          v-else
          v-field="form.jobs"
          :ref="setRef"
          :disabled="disabled"
          class="mt-3"
        />
      </template>

      <template #patterns="{ setRef }">
        <remediation-instruction-patterns-form
          v-field="form.patterns"
          :ref="setRef"
        />
      </template>
    </c-form-general-patterns-tabs>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { TEMPLATE_TESTING_TEST_TYPES, REMEDIATION_INSTRUCTION_TYPES } from '@/constants';

import { useI18n } from '@/hooks/i18n';

import RemediationInstructionGeneralForm from './remediation-instruction-general-form.vue';
import RemediationInstructionPatternsForm from './remediation-instruction-patterns-form.vue';
import RemediationInstructionStepsForm from './remediation-instruction-steps-form.vue';
import RemediationInstructionJobsForm from './remediation-instruction-jobs-form.vue';

export default {
  inject: ['$validator'],
  components: {
    RemediationInstructionGeneralForm,
    RemediationInstructionPatternsForm,
    RemediationInstructionStepsForm,
    RemediationInstructionJobsForm,
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
  setup(props) {
    const type = TEMPLATE_TESTING_TEST_TYPES.instruction;

    const { t, tc } = useI18n();

    const isManualType = computed(() => props.form.type === REMEDIATION_INSTRUCTION_TYPES.manual);

    const additionalLabel = computed(() => (isManualType.value ? tc('common.step', 2) : t('remediation.tabs.jobs')));

    return {
      type,

      isManualType,
      additionalLabel,
    };
  },
};
</script>
