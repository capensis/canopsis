<template>
  <div>
    <v-layout>
      <v-flex xs3>
        <c-instruction-type-field
          v-field="form.type"
          :disabled="disabled || !isNew"
          class="mb-2"
        />
      </v-flex>
      <v-flex>
        <c-enabled-field
          v-field="form.enabled"
          :disabled="disabledCommon"
          class="mt-0"
          hide-details
        />
      </v-flex>
    </v-layout>
    <c-name-field
      v-field="form.name"
      :disabled="disabledCommon"
    />
    <v-text-field
      v-field="form.description"
      v-validate="'required'"
      :label="$t('common.description')"
      :error-messages="errors.collect('description')"
      :disabled="disabledCommon"
      name="description"
    />
    <v-layout
      justify-space-between
      align-center
    >
      <v-flex xs7>
        <c-duration-field
          v-field="form.timeout_after_execution"
          :label="$t('remediation.instruction.timeoutAfterExecution')"
          :units-label="$t('common.unit')"
          :disabled="disabled"
          name="timeout_after_execution"
          required
        />
      </v-flex>
      <v-flex
        v-if="isAutoType"
        class="ml-2"
        xs3
      >
        <c-priority-field
          v-field="form.priority"
          :disabled="disabled"
        />
      </v-flex>
    </v-layout>
    <c-triggers-field
      v-if="isAutoType"
      v-field="form.triggers"
      :triggers="availableTriggers"
    />
    <remediation-instruction-jobs-form
      v-if="isAutoType || isManualSimplified"
      v-field="form.jobs"
      :disabled="disabled"
    />
    <remediation-instruction-steps-form
      v-else
      v-field="form.steps"
      :disabled="disabled"
    />
    <remediation-instruction-approval-form
      v-if="!disabledCommon"
      v-field="form.approval"
      :disabled="disabled"
      :required="requiredApprove"
    />
  </div>
</template>

<script>
import { REMEDIATION_AUTO_INSTRUCTION_TRIGGERS_TYPES } from '@/constants';

import { isInstructionAuto, isInstructionSimpleManual } from '@/helpers/entities/remediation/instruction/form';

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
    isNew: {
      type: Boolean,
      default: false,
    },
    requiredApprove: {
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
      return REMEDIATION_AUTO_INSTRUCTION_TRIGGERS_TYPES;
    },
  },
};
</script>
