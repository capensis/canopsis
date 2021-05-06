<template lang="pug">
  div
    v-layout(row)
      v-flex(xs3)
        c-instruction-type-field(v-field="form.type", @input="errors.clear()")
      v-flex
        c-enabled-field.mt-0(v-field="form.enabled", hide-details)
    v-layout(row)
      v-text-field(
        v-field="form.name",
        v-validate="'required'",
        :label="$t('common.name')",
        :error-messages="errors.collect('name')",
        name="name"
      )
    v-layout(row)
      v-text-field(
        v-field="form.description",
        v-validate="'required'",
        :label="$t('common.description')",
        :error-messages="errors.collect('description')",
        name="description"
      )
    v-layout(v-if="isManualType", row)
      remediation-instruction-steps-form(v-field="form.steps")
    template(v-else)
      v-layout(row, justify-space-between, align-center)
        v-flex(xs7)
          c-duration-field(
            v-field="form.timeout_after_execution",
            :label="$t('remediationInstructions.timeoutAfterExecution')",
            name="timeout_after_execution",
            required
          )
        v-flex.ml-2(xs3)
          c-priority-field(v-model="form.priority", required)
      v-layout(row)
        remediation-instruction-jobs-form(v-model="form.jobs")
    v-layout(row)
      remediation-instruction-approval-form(v-field="form.approval")
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
  },
  computed: {
    isManualType() {
      return this.form.type === REMEDIATION_INSTRUCTION_TYPES.manual;
    },
  },
};
</script>
