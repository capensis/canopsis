<template>
  <c-layout gap="2" column>
    <c-draggable-list-field
      :model-value="operations"
      :disabled="disabled"
      :class="{ 'grey lighten-1': isDragging }"
      class="flex-column gap-2"
      ghost-class="grey"
      handle=".operation-drag-handler"
      item-key="key"
      @update:model-value="updateModel"
      @start="startDragging"
      @end="endDragging"
    >
      <template #item="{ element: operation, index }">
        <remediation-instruction-operation-field
          :operation="operation"
          :operation-number="getOperationNumber(index)"
          :disabled="disabled"
          :template-vars="templateVars"
          @update:operation="updateItemInArray(index, $event)"
          @remove="removeOperation(index)"
        />
      </template>
    </c-draggable-list-field>
    <c-btn-with-error
      :error="hasOperationsErrors ? $t('remediation.instruction.errors.operationRequired') : ''"
      :disabled="disabled"
      variant="outlined"
      @click="addOperation"
    >
      {{ $t('remediation.instruction.addOperation') }}
    </c-btn-with-error>
  </c-layout>
</template>

<script>
import { computed, ref, watch } from 'vue';

import {
  remediationInstructionStepOperationToForm,
  getOperationNumber,
} from '@/helpers/entities/remediation/instruction/form';

import { useFormErrorsMap } from '@/hooks/form';
import { useArrayModelField } from '@/hooks/form/array-model-field';
import { useValidationAttachMinValueForField } from '@/hooks/validator/validation-attach-min-value';

import RemediationInstructionOperationField from './fields/remediation-instruction-operation-field.vue';

export default {
  components: {
    RemediationInstructionOperationField,
  },
  props: {
    name: {
      type: String,
      default: 'operations',
    },
    operations: {
      type: Array,
      default: () => [],
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
  emits: ['update:operations'],
  setup(props, { emit }) {
    const isDragging = ref(false);

    const { hasError } = useFormErrorsMap();
    const { updateModel, addItemIntoArray, updateItemInArray, removeItemFromArray } =
      useArrayModelField(props, emit);

    const hasOperationsErrors = computed(() => hasError(props.name));

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
      updateModel,
      updateItemInArray,
      addOperation,
      removeOperation: removeItemFromArray,
      startDragging,
      endDragging,
    };
  },
};
</script>
