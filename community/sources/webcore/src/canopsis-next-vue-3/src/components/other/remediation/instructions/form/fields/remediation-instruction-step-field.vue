<template>
  <c-card-iterator-item
    :item-number="stepNumber"
    class="remediation-instruction-step-field"
    @remove="remove"
  >
    <template #header="">
      <c-layout class="mt-3">
        <v-flex xs9>
          <v-text-field
            v-model="nameField"
            v-validate="'required'"
            :label="$t('common.name')"
            :error-messages="errors.collect(nameFieldName)"
            :name="nameFieldName"
            :disabled="disabled"
            filled
          />
        </v-flex>
        <v-flex class="pl-2" xs3>
          <v-text-field
            :model-value="formattedTimeToComplete"
            :label="$t('remediation.instruction.timeToComplete')"
            class="remediation-instruction-step-field__time-to-complete"
            readonly
          />
        </v-flex>
      </c-layout>
    </template>
    <c-form-block>
      <c-form-block-row :label="$t('common.workflow')" indented>
        <c-workflow-field
          v-model="stopOnFailField"
          :label="$t('remediation.instruction.workflow')"
          :continue-label="$t('remediation.instruction.remainingStep')"
          :disabled="disabled"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('remediation.instruction.endpoint')">
        <remediation-instruction-step-endpoint-field
          v-model="endpointField"
          :name="endpointFieldName"
          :disabled="disabled"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('remediation.instruction.operations')" indented>
        <remediation-instruction-operations-form
          :operations="step.operations"
          :name="operationFieldName"
          :step-number="stepNumber"
          :disabled="disabled"
          :template-vars="templateVars"
          @update:operations="updateStepOperations"
        />
      </c-form-block-row>
    </c-form-block>
  </c-card-iterator-item>
</template>

<script>
import { computed, toRef } from 'vue';

import { DATETIME_FORMATS } from '@/constants';
import { isOmitEqual } from '@/helpers/collection';
import { convertDurationToString, toSeconds } from '@/helpers/date/duration';
import { remediationInstructionStepToForm } from '@/helpers/entities/remediation/instruction/form';
import { useConfirmableForm } from '@/hooks/confirmable-form';

import RemediationInstructionOperationsForm from '../remediation-instruction-operations-form.vue';

import RemediationInstructionStepEndpointField from './remediation-instruction-step-endpoint-field.vue';

export default {
  inject: ['$validator'],
  components: {
    RemediationInstructionOperationsForm,
    RemediationInstructionStepEndpointField,
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
  emits: ['update:step', 'remove'],
  setup(props, { emit }) {
    const fieldSuffix = computed(() => (props.step.key ? `-${props.step.key}` : ''));
    const nameFieldName = computed(() => `name${fieldSuffix.value}`);
    const endpointFieldName = computed(() => `endpoint${fieldSuffix.value}`);
    const operationFieldName = computed(() => `operations${fieldSuffix.value}`);
    const timeToComplete = computed(() =>
      props.step.operations.reduce((acc, operation) => {
        const {
          time_to_complete: { value, unit },
        } = operation;

        return acc + toSeconds(value, unit);
      }, 0),
    );
    const createFieldModel = field =>
      computed({
        get: () => props.step[field],
        set: value =>
          emit('update:step', {
            ...props.step,
            [field]: value,
          }),
      });
    const nameField = createFieldModel('name');
    const stopOnFailField = createFieldModel('stop_on_fail');
    const endpointField = createFieldModel('endpoint');
    const formattedTimeToComplete = computed(() =>
      convertDurationToString(timeToComplete.value, DATETIME_FORMATS.refreshFieldFormat),
    );

    const { confirmAction: remove } = useConfirmableForm({
      form: toRef(props, 'step'),
      action: () => emit('remove'),
      comparator: (step) => {
        const emptyStep = remediationInstructionStepToForm();
        const paths = ['key', step.operations.length ? ['operations', 0, 'key'] : 'operations'];

        return isOmitEqual(step, emptyStep, paths);
      },
    });

    const updateStepOperations = operations =>
      emit('update:step', {
        ...props.step,
        operations,
      });

    return {
      nameFieldName,
      endpointFieldName,
      operationFieldName,
      nameField,
      stopOnFailField,
      endpointField,
      timeToComplete,
      formattedTimeToComplete,
      updateStepOperations,
      remove,
    };
  },
};
</script>

<style lang="scss">
.remediation-instruction-step-field {
  &__time-to-complete .v-input__slot {
    &:before,
    &:after {
      content: none !important;
    }
  }
}
</style>
