<template>
  <c-layout class="gap-3" column>
    <c-enabled-field v-model="enabledField" :disabled="disabledCommon" with-background />

    <c-form-general-patterns-tabs
      :form="formValues"
      :rule-id="ruleId"
      :type="type"
      :additional-label="additionalLabel"
    >
      <template #general="{ setRef, templateVars }">
        <remediation-instruction-general-form
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
          :ref="setRef"
          :disabled="disabled"
          :template-vars="templateVars"
          class="mt-3"
        />
        <remediation-instruction-jobs-form
          v-else
          :ref="setRef"
          :disabled="disabled"
          class="mt-3"
        />
      </template>

      <template #patterns="{ setRef }">
        <remediation-instruction-patterns-form :ref="setRef" />
      </template>
    </c-form-general-patterns-tabs>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { TEMPLATE_TESTING_TEST_TYPES, REMEDIATION_INSTRUCTION_TYPES } from '@/constants';

import { useFormContext, useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import RemediationInstructionGeneralForm from './remediation-instruction-general-form.vue';
import RemediationInstructionJobsForm from './remediation-instruction-jobs-form.vue';
import RemediationInstructionPatternsForm from './remediation-instruction-patterns-form.vue';
import RemediationInstructionStepsForm from './remediation-instruction-steps-form.vue';

export default {
  components: {
    RemediationInstructionGeneralForm,
    RemediationInstructionPatternsForm,
    RemediationInstructionStepsForm,
    RemediationInstructionJobsForm,
  },
  props: {
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
  setup() {
    const type = TEMPLATE_TESTING_TEST_TYPES.instruction;

    const { t } = useI18n();
    const form = useFormContext();
    const { enabledField, typeField } = useFormModels({
      enabledField: 'enabled',
      typeField: 'type',
    });

    const isManualType = computed(() => typeField.value === REMEDIATION_INSTRUCTION_TYPES.manual);
    const additionalLabel = computed(() => isManualType.value ? t('common.step', 2) : t('remediation.tabs.jobs'),
    );

    return {
      type,
      formValues: form.values,
      isManualType,
      enabledField,
      additionalLabel,
    };
  },
};
</script>
