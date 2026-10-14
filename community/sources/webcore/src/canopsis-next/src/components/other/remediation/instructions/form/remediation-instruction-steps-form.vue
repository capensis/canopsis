<template>
  <c-card-iterator-form
    :value="stepsField"
    :disabled="disabled"
    :draggable-group="draggableGroup"
    :required-error-message="t('remediation.instruction.errors.stepRequired')"
    :add-button-label="t('remediation.instruction.addStep')"
    :name="name"
    item-key="key"
    iterator-class="mb-2"
    required
    @add="addStep"
    @update:value="stepsField = $event"
  >
    <template #item="{ index, item: step }">
      <remediation-instruction-step-field
        :step="stepsField[index]"
        :name="`${name}.${step.key}`"
        :step-number="index + 1"
        :disabled="disabled"
        :template-vars="templateVars"
        @update:step="updateStep(index, $event)"
        @remove="removeStep(index)"
      />
    </template>
  </c-card-iterator-form>
</template>

<script>
import { remediationInstructionStepToForm } from '@/helpers/entities/remediation/instruction/form';

import { useFormArray } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import RemediationInstructionStepField from './fields/remediation-instruction-step-field.vue';

export default {
  components: {
    RemediationInstructionStepField,
  },
  props: {
    name: {
      type: String,
      default: 'steps',
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
  setup() {
    const { t } = useI18n();
    const stepsField = useFormArray('steps');

    const draggableGroup = {
      name: 'remediation-instruction-steps',
    };

    const addStep = () => {
      stepsField.value = [...stepsField.value, remediationInstructionStepToForm()];
    };
    const updateStep = (index, value) => {
      const steps = [...stepsField.value];

      steps[index] = value;
      stepsField.value = steps;
    };
    const removeStep = (index) => {
      stepsField.value = stepsField.value.filter((step, stepIndex) => stepIndex !== index);
    };

    return {
      t,
      stepsField,
      draggableGroup,
      updateStep,
      addStep,
      removeStep,
    };
  },
};
</script>
