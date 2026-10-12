<template>
  <c-lazy-search-field
    :value="fieldValue"
    :items="items"
    :label="$t('common.instructions')"
    :loading="pending"
    :has-more="hasMoreInstructions"
    item-text="name"
    item-value="_id"
    hide-details
    autocomplete
    no-auto-scroll
    @fetch="fetchInstructionsList"
    @fetch:more="fetchMoreInstructions"
    @update:search="updateInstructionsSearch"
    @update:value="updateValue"
    @keydown.capture="onStatisticInstructionBackspace"
    @focusout="restoreEmptyDefault"
  >
    <template #item="{ item, attrs, search }">
      <v-divider v-if="item.divider" />
      <v-list-item v-else class="c-lazy-search-field--tile" v-bind="attrs">
        <v-list-item-title>
          <v-list-item-mask :text="String(item.name ?? '')" :mask="search" />
        </v-list-item-title>
      </v-list-item>
    </template>
  </c-lazy-search-field>
</template>

<script setup>
import { keyBy, pick, uniq } from 'lodash';
import { computed, onMounted, ref, toRef, watch } from 'vue';

import { PAGINATION_LIMIT } from '@/config';
import { REMEDIATION_INSTRUCTION_TYPES } from '@/constants';

import { useLazySearchEmptyDefaultValue } from '@/hooks/form/lazy-search-empty-default-value';
import { useI18n } from '@/hooks/i18n';
import { useRemediationInstruction } from '@/hooks/store/modules/remediation-instruction';

const props = defineProps({
  modelValue: {
    type: [String, Number],
    required: true,
  },
});

const emit = defineEmits(['update:modelValue']);

const { t } = useI18n();
const { fetchRemediationInstructionsListWithoutStore } = useRemediationInstruction();
const {
  emptyValue,
  fieldValue,
  toFieldValue,
  updateValue,
  restoreEmptyDefault,
} = useLazySearchEmptyDefaultValue(
  {
    value: toRef(props, 'modelValue'),
  },
  emit,
);

const instructionsById = ref({});
const instructionIds = ref([]);
const pending = ref(false);
const page = ref(1);
const pageCount = ref(1);
const search = ref('');
let fetchRequestNumber = 0;

const statisticInstructionItems = computed(() => [
  {
    _id: emptyValue,
    name: t('remediation.statistic.allInstructions'),
  },
  {
    _id: REMEDIATION_INSTRUCTION_TYPES.manual,
    name: t('remediation.statistic.manualInstructions'),
  },
  {
    _id: REMEDIATION_INSTRUCTION_TYPES.auto,
    name: t('remediation.statistic.autoInstructions'),
  },
]);

/**
 * Checks whether the value is one of the statistic options (all / manual / auto).
 *
 * @param {string|number|symbol|null|undefined} instruction - Selected instruction id or statistic value.
 * @returns {boolean}
 */
const isStatisticInstructionItem = instruction =>
  statisticInstructionItems.value.some(({ _id }) => _id === toFieldValue(instruction));

/**
 * Checks whether the search text matches a statistic option label rather than a real query.
 *
 * @param {string} nextSearch - Current autocomplete search string.
 * @returns {boolean}
 */
const isStatisticInstructionSearch = nextSearch =>
  statisticInstructionItems.value.some(({ name }) => name === nextSearch);

/**
 * Clears the whole statistic option label on Backspace instead of deleting it character by character.
 *
 * @param {KeyboardEvent} event - Native keydown event from the autocomplete input.
 */
const onStatisticInstructionBackspace = (event) => {
  if (event.key !== 'Backspace') {
    return;
  }

  const { target } = event;

  if (target?.tagName !== 'INPUT' || !isStatisticInstructionSearch(target.value)) {
    return;
  }

  event.preventDefault();
  event.stopImmediatePropagation();

  target.value = '';
  target.dispatchEvent(new Event('input', { bubbles: true }));
};

const selectedInstructionIds = computed(() => {
  if (isStatisticInstructionItem(props.modelValue)) {
    return [];
  }

  return [props.modelValue];
});

const instructions = computed(() => instructionIds.value.map(id => instructionsById.value[id]).filter(Boolean),
);

/**
 * Fetches the currently selected instruction when it is missing from the loaded list.
 *
 * @returns {Promise<void>}
 */
const fetchSelectedInstructions = async () => {
  const missingSelectedIds = selectedInstructionIds.value.filter(id => !instructionsById.value[id]);

  if (!missingSelectedIds.length) {
    return;
  }

  const { data = [] } = await fetchRemediationInstructionsListWithoutStore({
    params: {
      limit: missingSelectedIds.length,
      ids: missingSelectedIds,
    },
  });

  instructionsById.value = {
    ...instructionsById.value,
    ...keyBy(data, '_id'),
  };
  instructionIds.value = uniq([...instructionIds.value, ...selectedInstructionIds.value]);
};

/**
 * Fetches a page of instructions. Statistic option labels are not sent as `search`.
 *
 * @param {Object} [options]
 * @param {number} [options.nextPage=1] - Page number to request.
 * @param {string} [options.search] - Search query. Ignored when it matches a statistic option label.
 * @returns {Promise<void>}
 */
const fetchInstructionsList = async ({
  nextPage = 1,
  search: nextSearch = search.value,
} = {}) => {
  const requestNumber = ++fetchRequestNumber;

  try {
    pending.value = true;

    const params = {
      limit: PAGINATION_LIMIT,
      page: nextPage,
    };

    if (nextSearch && !isStatisticInstructionSearch(nextSearch)) {
      params.search = nextSearch;
    }

    const { data = [], meta = {} } = await fetchRemediationInstructionsListWithoutStore({ params });

    if (requestNumber !== fetchRequestNumber) {
      return;
    }

    page.value = nextPage;
    pageCount.value = meta.page_count ?? 1;
    search.value = nextSearch;
    instructionsById.value = {
      ...(nextPage !== 1 ? instructionsById.value : pick(instructionsById.value, selectedInstructionIds.value)),
      ...keyBy(data, '_id'),
    };
    instructionIds.value = uniq([
      ...(nextPage !== 1 ? instructionIds.value : []),
      ...data.map(instruction => instruction._id),
      ...selectedInstructionIds.value,
    ]);

    await fetchSelectedInstructions();
  } catch (err) {
    console.error(err);
  } finally {
    if (requestNumber === fetchRequestNumber) {
      pending.value = false;
    }
  }
};

/**
 * Loads the next page of instructions when more results are available.
 *
 * @returns {Promise<void>|undefined}
 */
const fetchMoreInstructions = () => {
  if (hasMoreInstructions.value) {
    return fetchInstructionsList({ nextPage: page.value + 1 });
  }

  return undefined;
};

/**
 * Reloads the first page of instructions for a new search. Statistic option labels are not queried.
 *
 * @param {string} nextSearch - Search string emitted by the autocomplete.
 * @returns {Promise<void>|undefined}
 */
const updateInstructionsSearch = (nextSearch) => {
  if (isStatisticInstructionSearch(nextSearch)) {
    if (!search.value) {
      return undefined;
    }

    return fetchInstructionsList({
      nextPage: 1,
      search: '',
    });
  }

  return fetchInstructionsList({
    nextPage: 1,
    search: nextSearch ?? '',
  });
};

const hasMoreInstructions = computed(() => pageCount.value > page.value);

const items = computed(() => [
  ...statisticInstructionItems.value,
  {
    _id: '__divider__',
    divider: true,
  },
  ...instructions.value,
]);

watch(() => props.modelValue, fetchSelectedInstructions);

onMounted(fetchInstructionsList);
</script>
