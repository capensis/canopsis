import { isNil, keyBy, pick, uniq } from 'lodash';
import { ref, unref, computed, onMounted, watch } from 'vue';

import { PAGINATION_LIMIT } from '@/config';
import {
  REMEDIATION_INSTRUCTION_FILTER_FIELDS,
  REMEDIATION_INSTRUCTION_FILTER_TYPES,
} from '@/constants';

import {
  isAutoInstructionType,
  isManualInstructionType,
  isAnyManualInstructionType,
} from '@/helpers/entities/remediation/instruction/list';

import { useRemediationInstruction } from '@/hooks/store/modules/remediation-instruction';

/**
 * Provides logic to fetch and filter remediation instructions for alarms list.
 *
 * @param {Ref|Object} filter - Reactive filter object for instructions.
 * @returns {{
 *   instructions: ComputedRef<Array>,
 *   instructionsPending: Ref<boolean>,
 *   initialInstructionsPending: Ref<boolean>,
 *   hasMoreInstructions: ComputedRef<boolean>,
 *   fetchInstructionsList: Function,
 *   fetchMoreInstructions: Function,
 *   updateInstructionsSearch: Function
 * }}
 */
export const useAlarmsListRemediationInstructionsFilterFetch = (filter) => {
  const instructionsById = ref({});
  const instructionsIds = ref([]);
  const instructionsPending = ref(false);
  const initialInstructionsPending = ref(false);
  const page = ref(1);
  const pageCount = ref(1);
  const search = ref('');
  let fetchRequestNumber = 0;

  const { fetchRemediationInstructionsListWithoutStore } = useRemediationInstruction();

  const allInstructions = computed(() => instructionsIds.value.map(id => instructionsById.value[id]).filter(Boolean),
  );

  const getSelectedInstructionIds = () => unref(filter).instruction_ids ?? [];

  const fetchSelectedInstructions = async () => {
    const selectedIds = getSelectedInstructionIds();

    if (!selectedIds.length) {
      return;
    }

    const missingSelectedIds = selectedIds.filter(id => !instructionsById.value[id]);

    if (!missingSelectedIds.length) {
      return;
    }

    const { data } = await fetchRemediationInstructionsListWithoutStore({
      params: {
        limit: missingSelectedIds.length,
        ids: missingSelectedIds,
      },
    });

    instructionsById.value = {
      ...instructionsById.value,
      ...keyBy(data, '_id'),
    };
    instructionsIds.value = uniq([...instructionsIds.value, ...selectedIds]);
  };

  const fetchInstructionsList = async ({
    nextPage = 1,
    search: nextSearch = search.value,
  } = {}) => {
    const requestNumber = ++fetchRequestNumber;

    try {
      instructionsPending.value = true;

      const { data, meta } = await fetchRemediationInstructionsListWithoutStore({
        params: {
          limit: PAGINATION_LIMIT,
          page: nextPage,
          search: nextSearch,
        },
      });

      if (requestNumber !== fetchRequestNumber) {
        return;
      }

      const selectedIds = getSelectedInstructionIds();

      page.value = nextPage;
      pageCount.value = meta?.page_count ?? 1;
      search.value = nextSearch;
      instructionsById.value = {
        ...(nextPage !== 1 ? instructionsById.value : pick(instructionsById.value, selectedIds)),
        ...keyBy(data, '_id'),
      };
      instructionsIds.value = uniq([
        ...(nextPage !== 1 ? instructionsIds.value : []),
        ...data.map(instruction => instruction._id),
        ...selectedIds,
      ]);

      await fetchSelectedInstructions();
    } catch (err) {
      console.error(err);
    } finally {
      if (requestNumber === fetchRequestNumber) {
        instructionsPending.value = false;
      }
    }
  };

  const fetchMoreInstructions = () => {
    if (pageCount.value > page.value) {
      return fetchInstructionsList({ nextPage: page.value + 1 });
    }

    return undefined;
  };

  const updateInstructionsSearch = nextSearch => fetchInstructionsList({
      nextPage: 1,
      search: nextSearch ?? '',
    });

  const preparedInstructions = computed(() => {
    const unwrappedFilter = unref(filter);

    const isManualInstructionSelectedType = isManualInstructionType(
      unwrappedFilter.instruction_type,
    );
    const isAutoInstructionSelectedType = isAutoInstructionType(unwrappedFilter.instruction_type);

    return allInstructions.value.filter(
      instruction => (isManualInstructionSelectedType && isAnyManualInstructionType(instruction.type)) ||
        (isAutoInstructionSelectedType && isAutoInstructionType(instruction.type)),
    );
  });

  const hasMoreInstructions = computed(() => pageCount.value > page.value);

  watch(() => unref(filter).instruction_ids, fetchSelectedInstructions);

  watch(
    () => unref(filter).instruction_type,
    () => fetchInstructionsList({ nextPage: 1 }),
  );

  onMounted(async () => {
    initialInstructionsPending.value = true;

    await fetchSelectedInstructions();

    initialInstructionsPending.value = false;

    await fetchInstructionsList();
  });

  return {
    instructions: preparedInstructions,
    instructionsPending,
    initialInstructionsPending,
    hasMoreInstructions,
    fetchInstructionsList,
    fetchMoreInstructions,
    updateInstructionsSearch,
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

    if (
      unwrappedFilter.instruction_filter_type !==
      REMEDIATION_INSTRUCTION_FILTER_TYPES.hasInstructions
    ) {
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
