<template>
  <v-layout class="gap-2" column>
    <c-draggable-list-field
      v-field="operations"
      :disabled="disabled"
      :class="{ 'grey lighten-1': isDragging }"
      class="flex-column gap-2"
      ghost-class="grey"
      handle=".operation-drag-handler"
      @start="startDragging"
      @end="endDragging"
    >
      <remediation-instruction-operation-field
        v-for="(operation, index) in operations"
        v-field="operations[index]"
        :key="operation.key"
        :index="index"
        :operation-number="getOperationNumber(index)"
        :disabled="disabled"
        :template-vars="templateVars"
        @remove="removeOperation(index)"
      />
    </c-draggable-list-field>
    <c-btn-with-error
      :error="hasOperationsErrors ? $t('remediation.instruction.errors.operationRequired') : ''"
      :disabled="disabled"
      outlined
      @click="addOperation"
    >
      {{ $t('remediation.instruction.addOperation') }}
    </c-btn-with-error>
  </v-layout>
</template>

<script>
import { computed, ref, watch } from 'vue';

import {
  remediationInstructionStepOperationToForm,
  getOperationNumber,
} from '@/helpers/entities/remediation/instruction/form';

import { useArrayModelField } from '@/hooks/form/array-model-field';
import { useValidator } from '@/hooks/validator/validator';
import { useValidationAttachMinValueForField } from '@/hooks/validator/validation-attach-min-value';

import RemediationInstructionOperationField from './fields/remediation-instruction-operation-field.vue';

export default {
  components: {
    RemediationInstructionOperationField,
  },
  model: {
    prop: 'operations',
    event: 'input',
  },
  props: {
    name: {
      type: String,
      default: 'operations',
    },
    operations: {
      type: Array,
      default: () => ([]),
    },
    stepNumber: {
      type: [String, Number],
      default: '',
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
    const isDragging = ref(false);

    const validator = useValidator();
    const { addItemIntoArray, removeItemFromArray } = useArrayModelField(props, emit);

    const hasOperationsErrors = computed(() => validator?.errors?.has(props.name) ?? false);

    const { asyncValidateMinValueRule } = useValidationAttachMinValueForField(
      props.name,
      () => props.operations.length,
    );

    watch(() => props.operations, asyncValidateMinValueRule);

    /**
     * Appends a new empty operation to the operations list.
     */
    const addOperation = () => addItemIntoArray(remediationInstructionStepOperationToForm());

    /**
     * Returns the display number for an operation within the current step.
     *
     * @param {number} index - Zero-based index of the operation in the list.
     * @returns {string}
     */
    const getOperationNumberByIndex = index => getOperationNumber(props.stepNumber, index);

    /**
     * Marks the draggable list as being dragged to apply visual feedback.
     */
    const startDragging = () => {
      isDragging.value = true;
    };

    /**
     * Clears the dragging state when the drag operation ends.
     */
    const endDragging = () => {
      isDragging.value = false;
    };

    return {
      isDragging,
      hasOperationsErrors,
      getOperationNumber: getOperationNumberByIndex,
      addOperation,
      removeOperation: removeItemFromArray,
      startDragging,
      endDragging,
    };
  },
};
</script>
