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
    v-layout(row)
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
          :label="$t('remediationInstructions.timeoutAfterExecution')",
          :disabled="disabled",
          name="timeout_after_execution",
          required
        )
      v-flex.ml-2(v-if="!isManualType", xs3)
        c-priority-field(v-model="form.priority", :disabled="disabled", required)
    v-layout(v-if="isManualType", row)
      remediation-instruction-steps-form(v-field="form.steps", :disabled="disabled")
    v-layout(v-else, row)
      remediation-instruction-jobs-form(v-model="form.jobs", :disabled="disabled")
    v-layout(v-if="!disabledCommon", row)
      remediation-instruction-approval-form(v-field="form.approval", :disabled="disabled")
</template>

<script>
import { REMEDIATION_INSTRUCTION_TYPES } from '@/constants';

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
    isManualType() {
      return this.form.type === REMEDIATION_INSTRUCTION_TYPES.manual;
    },
  },
};
</script>
