import { isNil } from 'lodash';
import { ref, unref, computed, onMounted } from 'vue';

import { MAX_LIMIT, REMEDIATION_INSTRUCTION_FILTER_FIELDS, REMEDIATION_INSTRUCTION_FILTER_TYPES } from '@/constants';

import {
  isAutoInstructionType,
  isManualInstructionType,
  isAnyManualInstructionType,
} from '@/helpers/entities/remediation/instruction/list';

import { usePendingHandler } from '@/hooks/query/pending';
import { useRemdeitionInstruction } from '@/hooks/store/modules/remediation-instruction';

/**
 * Provides logic to fetch and filter remediation instructions for alarms list.
 *
 * @param {Ref|Object} filter - Reactive filter object for instructions.
 * @returns {{
 *   instructions: ComputedRef<Array>,
 *   instructionsPending: Ref<boolean>
 * }}
 */
export const useAlarmsListRemediationInstructionsFilterFetch = (filter) => {
  const instructions = ref([]);

  const { fetchRemediationInstructionsListWithoutStore } = useRemdeitionInstruction();

  const { handler: fetchInstructionsList, pending: instructionsPending } = usePendingHandler(async () => {
    const { data } = await fetchRemediationInstructionsListWithoutStore({ params: { limit: MAX_LIMIT } });

    instructions.value = data;
  });

  const preparedInstructions = computed(() => {
    const unwrappedFilter = unref(filter);

    const isManualInstructionSelectedType = isManualInstructionType(unwrappedFilter.instruction_type);
    const isAutoInstructionSelectedType = isAutoInstructionType(unwrappedFilter.instruction_type);

    return instructions.value
      .filter(instruction => (
        isManualInstructionSelectedType && isAnyManualInstructionType(instruction.type)
      ) || (
        isAutoInstructionSelectedType && isAutoInstructionType(instruction.type)
      ));
  });

  onMounted(fetchInstructionsList);

  return {
    instructions: preparedInstructions,
    instructionsPending,
  };
};

/**
 * Provides computed fields and selection logic for remediation instruction filters in alarms list.
 *
 * @param {Ref|Object} filter - Reactive filter object for instructions.
 * @returns {{
 *   fields: ComputedRef<Array>,
 *   selectedItems: ComputedRef<Array>,
 *   selectedAllFields: ComputedRef<boolean>
 * }}
 */
export const useAlarmsListRemediationInstructionsFilterFields = (filter) => {
  const fields = computed(() => {
    const unwrappedFilter = unref(filter);
    const result = [REMEDIATION_INSTRUCTION_FILTER_FIELDS.instructionFilterType];

    if (unwrappedFilter.instruction_filter_type !== REMEDIATION_INSTRUCTION_FILTER_TYPES.hasInstructions) {
      return result;
    }

    if (!isNil(unwrappedFilter.instruction_filter_type)) {
      result.push(REMEDIATION_INSTRUCTION_FILTER_FIELDS.instructionType);

      if (!isNil(unwrappedFilter.instruction_type)) {
        result.push(
          REMEDIATION_INSTRUCTION_FILTER_FIELDS.instructionStatuses,
          REMEDIATION_INSTRUCTION_FILTER_FIELDS.instructionIds,
        );
      }
    }

    return result;
  });

  const selectedItems = computed(() => {
    const unwrappedFilter = unref(filter);
    const result = [];

    if (isNil(unwrappedFilter.instruction_filter_type)) {
      return result;
    }

    result.push(REMEDIATION_INSTRUCTION_FILTER_FIELDS.instructionFilterType);

    if (isNil(unwrappedFilter.instruction_type)) {
      return result;
    }

    result.push(REMEDIATION_INSTRUCTION_FILTER_FIELDS.instructionType);

    if (unwrappedFilter.instruction_statuses?.length) {
      result.push(REMEDIATION_INSTRUCTION_FILTER_FIELDS.instructionStatuses);
    }

    if (unwrappedFilter.instruction_ids?.length) {
      result.push(REMEDIATION_INSTRUCTION_FILTER_FIELDS.instructionIds);
    }

    return result;
  });

  const selectedAllFields = computed(() => selectedItems.value.length === fields.value.length);

  return {
    fields,
    selectedItems,
    selectedAllFields,
  };
};
