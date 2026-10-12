<template>
  <c-card-iterator-form
    v-field="steps"
    :disabled="disabled"
    :draggable-group="draggableGroup"
    :required-error-message="$t('remediation.instruction.errors.stepRequired')"
    :add-button-label="$t('remediation.instruction.addStep')"
    :name="name"
    item-key="key"
    iterator-class="mb-2"
    required
    @add="addStep"
  >
    <template #item="{ index }">
      <remediation-instruction-step-field
        v-field="steps[index]"
        :step-number="index + 1"
        :disabled="disabled"
        :template-vars="templateVars"
        @remove="removeStep(index)"
      />
    </template>
  </c-card-iterator-form>
</template>

<script>
import { remediationInstructionStepToForm } from '@/helpers/entities/remediation/instruction/form';

import { useArrayModelField } from '@/hooks/form/array-model-field';

import RemediationInstructionStepField from './fields/remediation-instruction-step-field.vue';

export default {
  components: {
    RemediationInstructionStepField,
  },
  model: {
    prop: 'steps',
    event: 'input',
  },
  props: {
    steps: {
      type: Array,
      default: () => ([]),
    },
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
  setup(props, { emit }) {
    const { addItemIntoArray, removeItemFromArray } = useArrayModelField(props, emit);

    const draggableGroup = {
      name: 'remediation-instruction-steps',
    };

    const addStep = () => addItemIntoArray(remediationInstructionStepToForm());

    return {
      draggableGroup,
      addStep,
      removeStep: removeItemFromArray,
    };
  },
};
</script>
