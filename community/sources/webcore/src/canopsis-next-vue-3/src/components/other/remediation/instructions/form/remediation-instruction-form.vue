<template>
  <c-layout class="gap-3" column>
    <c-enabled-field v-model="enabledField" :disabled="disabledCommon" with-background />

    <c-form-general-patterns-tabs
      :form="form"
      :rule-id="ruleId"
      :type="type"
      :additional-label="additionalLabel"
      :with-template-testing="false"
    >
      <template #general="{ setRef, templateVars }">
        <remediation-instruction-general-form
          :ref="setRef"
          :form="form"
          :disabled="disabled"
          :disabled-common="disabledCommon"
          :is-new="isNew"
          :required-approve="requiredApprove"
          :template-vars="templateVars"
          @update:form="updateModel"
        />
      </template>

      <template #additional="{ setRef, templateVars }">
        <remediation-instruction-steps-form
          v-if="isManualType"
          :ref="setRef"
          :steps="form.steps"
          :disabled="disabled"
          :template-vars="templateVars"
          class="mt-3"
          @update:steps="updateField('steps', $event)"
        />
        <remediation-instruction-jobs-form
          v-else
          :ref="setRef"
          :jobs="form.jobs"
          :disabled="disabled"
          class="mt-3"
          @update:jobs="updateField('jobs', $event)"
        />
      </template>

      <template #patterns="{ setRef }">
        <remediation-instruction-patterns-form
          :ref="setRef"
          :form="form.patterns"
          @update:form="updateField('patterns', $event)"
        />
      </template>
    </c-form-general-patterns-tabs>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { TEMPLATE_TESTING_TEST_TYPES, REMEDIATION_INSTRUCTION_TYPES } from '@/constants';

import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';

import RemediationInstructionGeneralForm from './remediation-instruction-general-form.vue';
import RemediationInstructionJobsForm from './remediation-instruction-jobs-form.vue';
import RemediationInstructionPatternsForm from './remediation-instruction-patterns-form.vue';
import RemediationInstructionStepsForm from './remediation-instruction-steps-form.vue';

export default {
  inject: ['$validator'],
  components: {
    RemediationInstructionGeneralForm,
    RemediationInstructionPatternsForm,
    RemediationInstructionStepsForm,
    RemediationInstructionJobsForm,
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
  emits: ['update:form'],
  setup(props, { emit }) {
    const type = TEMPLATE_TESTING_TEST_TYPES.instruction;

    const { t } = useI18n();
    const { updateModel, updateField } = useModelField(props, emit);

    const isManualType = computed(() => props.form.type === REMEDIATION_INSTRUCTION_TYPES.manual);
    const enabledField = computed({
      get: () => props.form.enabled,
      set: value => updateField('enabled', value),
    });

    const additionalLabel = computed(() => isManualType.value ? t('common.step', 2) : t('remediation.tabs.jobs'),
    );

    return {
      type,

      isManualType,
      enabledField,
      additionalLabel,
      updateModel,
      updateField,
    };
  },
};
</script>
