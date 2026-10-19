<template>
  <c-card-iterator-item
    :item-number="stepNumber"
    class="remediation-instruction-step-field"
    @remove="remove"
  >
    <template #header="">
      <v-layout class="mt-3">
        <v-flex xs9>
          <v-text-field
            v-field="step.name"
            v-validate="'required'"
            :label="$t('common.name')"
            :error-messages="errors.collect(nameFieldName)"
            :name="nameFieldName"
            :disabled="disabled"
            filled
          />
        </v-flex>
        <v-flex
          class="pl-2"
          xs3
        >
          <v-text-field
            :value="timeToComplete | duration('refreshFieldFormat')"
            :label="$t('remediation.instruction.timeToComplete')"
            class="remediation-instruction-step-field__time-to-complete"
            readonly
          />
        </v-flex>
      </v-layout>
    </template>
    <c-form-block>
      <c-form-block-row :label="$t('common.workflow')" indented>
        <c-workflow-field
          v-field="step.stop_on_fail"
          :label="$t('remediation.instruction.workflow')"
          :continue-label="$t('remediation.instruction.remainingStep')"
          :disabled="disabled"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('remediation.instruction.endpoint')">
        <remediation-instruction-step-endpoint-field
          v-field="step.endpoint"
          :name="endpointFieldName"
          :disabled="disabled"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('remediation.instruction.operations')" indented>
        <remediation-instruction-operations-form
          v-field="step.operations"
          :name="operationFieldName"
          :step-number="stepNumber"
          :disabled="disabled"
          :template-vars="templateVars"
        />
      </c-form-block-row>
    </c-form-block>
  </c-card-iterator-item>
</template>

<script>
import { computed, toRef } from 'vue';

import { remediationInstructionStepToForm } from '@/helpers/entities/remediation/instruction/form';
import { isOmitEqual } from '@/helpers/collection';
import { toSeconds } from '@/helpers/date/duration';

import { useConfirmableForm } from '@/hooks/confirmable-form';

import RemediationInstructionOperationsForm from '../remediation-instruction-operations-form.vue';

import RemediationInstructionStepEndpointField from './remediation-instruction-step-endpoint-field.vue';

export default {
  inject: ['$validator'],
  components: {
    RemediationInstructionOperationsForm,
    RemediationInstructionStepEndpointField,
  },
  model: {
    prop: 'step',
    event: 'input',
  },
  props: {
    step: {
      type: Object,
      required: true,
    },
    stepNumber: {
      type: [Number, String],
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    templateVars: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props, { emit }) {
    const fieldSuffix = computed(() => (props.step.key ? `-${props.step.key}` : ''));
    const nameFieldName = computed(() => `name${fieldSuffix.value}`);
    const endpointFieldName = computed(() => `endpoint${fieldSuffix.value}`);
    const operationFieldName = computed(() => `operations${fieldSuffix.value}`);
    const timeToComplete = computed(() => props.step.operations.reduce((acc, operation) => {
      const { time_to_complete: { value, unit } } = operation;

      return acc + toSeconds(value, unit);
    }, 0));

    const { confirmAction: remove } = useConfirmableForm({
      form: toRef(props, 'step'),
      action: () => emit('remove'),
      comparator: (step) => {
        const emptyStep = remediationInstructionStepToForm();
        const paths = [
          'key',
          step.operations.length ? ['operations', 0, 'key'] : 'operations',
        ];

        return isOmitEqual(step, emptyStep, paths);
      },
    });

    return {
      nameFieldName,
      endpointFieldName,
      operationFieldName,
      timeToComplete,
      remove,
    };
  },
};
</script>

<style lang="scss">
.remediation-instruction-step-field {
  &__time-to-complete .v-input__slot {
    &:before, &:after {
      content: none !important;
    }
  }
}
</style>
