<template>
  <c-lazy-search-field
    :value="modelValue"
    :items="items"
    :label="$t('common.instructions')"
    :loading="pending"
    :has-more="hasMoreInstructions"
    item-text="name"
    item-value="_id"
    hide-details
    autocomplete
    @fetch="fetchInstructionsList"
    @fetch:more="fetchMoreInstructions"
    @update:search="updateInstructionsSearch"
    @update:value="$emit('update:modelValue', $event)"
  >
    <template #item="{ item, attrs }">
      <v-divider v-if="item.divider" />
      <v-list-item v-else v-bind="attrs" />
    </template>
  </c-lazy-search-field>
</template>

<script setup>
import { keyBy, pick, uniq } from 'lodash';
import { computed, onMounted, ref, watch } from 'vue';

import { PAGINATION_LIMIT } from '@/config';
import { REMEDIATION_INSTRUCTION_TYPES } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useRemediationInstruction } from '@/hooks/store/modules/remediation-instruction';

const props = defineProps({
  modelValue: {
    type: [String, Number],
    required: true,
  },
});

defineEmits(['update:modelValue']);

const { t } = useI18n();
const { fetchRemediationInstructionsListWithoutStore } = useRemediationInstruction();

const instructionsById = ref({});
const instructionIds = ref([]);
const pending = ref(false);
const page = ref(1);
const pageCount = ref(1);
const search = ref('');
let fetchRequestNumber = 0;

const statisticInstructionItems = computed(() => [
  {
    _id: '',
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

const isStatisticInstructionItem = instruction => statisticInstructionItems.value.some(({ _id }) => _id === instruction);

const selectedInstructionIds = computed(() => {
  if (!props.modelValue || isStatisticInstructionItem(props.modelValue)) {
    return [];
  }

  return [props.modelValue];
});

const instructions = computed(() => instructionIds.value.map(id => instructionsById.value[id]).filter(Boolean),
);

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

const fetchInstructionsList = async ({
  nextPage = 1,
  search: nextSearch = search.value,
} = {}) => {
  const requestNumber = ++fetchRequestNumber;

  try {
    pending.value = true;

    const { data = [], meta = {} } = await fetchRemediationInstructionsListWithoutStore({
      params: {
        limit: PAGINATION_LIMIT,
        page: nextPage,
        search: nextSearch,
      },
    });

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

const fetchMoreInstructions = () => {
  if (hasMoreInstructions.value) {
    return fetchInstructionsList({ nextPage: page.value + 1 });
  }

  return undefined;
};

const updateInstructionsSearch = nextSearch => fetchInstructionsList({
    nextPage: 1,
    search: nextSearch ?? '',
  });

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
