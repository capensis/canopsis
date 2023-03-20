<template lang="pug">
  div
    v-layout(row)
      v-flex(xs3)
        c-instruction-type-field.mb-2(
          v-field="form.type",
          :disabled="disabled"
        )
      v-flex
        c-enabled-field.mt-0(
          v-field="form.enabled",
          :disabled="disabledCommon",
          hide-details
        )
    c-name-field(v-field="form.name", :disabled="disabledCommon")
    v-text-field(
      v-field="form.description",
      v-validate="'required'",
      :label="$t('common.description')",
      :error-messages="errors.collect('description')",
      :disabled="disabledCommon",
      name="description"
    )
    v-layout(row, justify-space-between, align-center)
      v-flex(xs7)
        c-duration-field(
          v-field="form.timeout_after_execution",
          :label="$t('remediation.instruction.timeoutAfterExecution')",
          :units-label="$t('common.unit')",
          :disabled="disabled",
          name="timeout_after_execution",
          required
        )
      v-flex.ml-2(v-if="isAutoType", xs3)
        c-priority-field(v-field="form.priority", :disabled="disabled", required)
    c-triggers-field(
      v-if="isAutoType",
      v-field="form.triggers",
      :triggers="availableTriggers"
    )
    remediation-instruction-jobs-form(
      v-if="isAutoType || isManualSimplified",
      v-field="form.jobs",
      :disabled="disabled"
    )
    remediation-instruction-steps-form(v-else, v-field="form.steps", :disabled="disabled")
    remediation-instruction-approval-form(v-if="!disabledCommon", v-field="form.approval", :disabled="disabled")
</template>

<script>
import { REMEDIATION_AUTO_INSTRUCTION_TRIGGERS } from '@/constants';

import { isInstructionAuto, isInstructionSimpleManual } from '@/helpers/forms/remediation-instruction';

import RemediationInstructionStepsForm from './remediation-instruction-steps-form.vue';
import RemediationInstructionJobsForm from './remediation-instruction-jobs-form.vue';
import RemediationInstructionApprovalForm from './remediation-instruction-approval-form.vue';

export default {
  inject: ['$validator'],
  components: {
    RemediationInstructionStepsForm,
    RemediationInstructionJobsForm,
    RemediationInstructionApprovalForm,
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
  },
  computed: {
    isAutoType() {
      return isInstructionAuto(this.form);
    },

    isManualSimplified() {
      return isInstructionSimpleManual(this.form);
    },

    availableTriggers() {
      return REMEDIATION_AUTO_INSTRUCTION_TRIGGERS;
    },
  },
};
</script>
