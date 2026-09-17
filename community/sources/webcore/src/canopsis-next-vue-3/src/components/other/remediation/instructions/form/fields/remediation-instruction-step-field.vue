<template>
  <c-card-iterator-item
    :item-number="stepNumber"
    class="remediation-instruction-step-field"
    static-delete
    @remove="remove"
  >
    <template #header>
      <v-row align="center">
        <v-col>
          <c-name-field
            :model-value="step.name"
            :name="nameFieldName"
            :disabled="disabled"
            variant="filled"
            density="comfortable"
            required
            @update:model-value="updateField('name', $event)"
          />
        </v-col>
        <v-col cols="3">
          <c-text-field
            :model-value="formattedTimeToComplete"
            :label="t('remediation.instruction.timeToComplete')"
            variant="plain"
            readonly
          />
        </v-col>
      </v-row>
    </template>
    <c-form-block>
      <c-form-block-row :label="t('common.workflow')">
        <c-workflow-field
          :model-value="step.stop_on_fail"
          :label="t('remediation.instruction.workflow')"
          :continue-label="t('remediation.instruction.remainingStep')"
          :disabled="disabled"
          @update:model-value="updateField('stop_on_fail', $event)"
        />
      </c-form-block-row>

      <c-form-block-row :label="t('remediation.instruction.endpoint')">
        <remediation-instruction-step-endpoint-field
          :model-value="step.endpoint"
          :name="endpointFieldName"
          :disabled="disabled"
          @update:model-value="updateField('endpoint', $event)"
        />
      </c-form-block-row>

      <c-form-block-row :label="t('remediation.instruction.operations')" indented>
        <remediation-instruction-operations-form
          :operations="step.operations"
          :name="operationFieldName"
          :step-number="stepNumber"
          :disabled="disabled"
          :template-vars="templateVars"
          @update:operations="updateField('operations', $event)"
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
import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';

import RemediationInstructionOperationsForm from '../remediation-instruction-operations-form.vue';

import RemediationInstructionStepEndpointField from './remediation-instruction-step-endpoint-field.vue';

export default {
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
    const { t } = useI18n();
    const { updateField } = useModelField(props, emit);

    const fieldSuffix = computed(() => (props.step.key ? `-${props.step.key}` : ''));
    const nameFieldName = computed(() => `name${fieldSuffix.value}`);
    const endpointFieldName = computed(() => `endpoint${fieldSuffix.value}`);
    const operationFieldName = computed(() => `operations${fieldSuffix.value}`);
    const timeToComplete = computed(() => props.step.operations.reduce((acc, operation) => {
        const {
          time_to_complete: { value, unit },
        } = operation;

        return acc + toSeconds(value, unit);
      }, 0),
    );
    const formattedTimeToComplete = computed(
      () => convertDurationToString(timeToComplete.value, DATETIME_FORMATS.refreshFieldFormat),
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

    return {
      t,
      nameFieldName,
      endpointFieldName,
      operationFieldName,
      formattedTimeToComplete,
      updateField,
      remove,
    };
  },
};
</script>

<style lang="scss" scoped>
.remediation-instruction-step-field {
  :deep(.c-card-iterator-item__row) {
    align-items: center;
  }
}
</style>
