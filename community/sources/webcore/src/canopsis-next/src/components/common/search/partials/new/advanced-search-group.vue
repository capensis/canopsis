<template>
  <span class="c-new-advanced-search__group">
    <v-chip
  </span>
</template>

<script>
import { isUndefined } from 'lodash';
import { computed, ref } from 'vue';

import AdvancedSearchChip from './advanced-search-chip.vue';

const ADVANCED_SEARCH_CHIP_TYPES = {
  field: 'field',
  operator: 'operator',
  valueType: 'valueType',
  value: 'value',
  valueSecond: 'valueSecond',
  union: 'union',
};

const example = {
  key: 'asdasdasd',
  field: {
    operators: [],
    fetch: () => {},
  },
  operator: '',
  valueType: '',
  value: '',
  secondValue: '',
};

const ADVANCED_SEARCH_CHIPS_GROUP_TYPES_ORDER = [
  ADVANCED_SEARCH_CHIP_TYPES.field,
  ADVANCED_SEARCH_CHIP_TYPES.valueType,
  ADVANCED_SEARCH_CHIP_TYPES.operator,
  ADVANCED_SEARCH_CHIP_TYPES.value,
  ADVANCED_SEARCH_CHIP_TYPES.valueSecond,
];

export default {
  components: { AdvancedSearchChip },
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props, { emit }) {
    const activeType = ref();
    const chips = computed(() => (
      ADVANCED_SEARCH_CHIPS_GROUP_TYPES_ORDER.map(type => props.value[type]).filter(chip => !isUndefined(chip))
    ));
    const hasValueType = false;
    const hasSecondValue = false;

    const isLastChip = index => index === chips.value.length - 1;
    const remove = () => emit('remove');

    return {
      chips,
      hasValueType,
      hasSecondValue,
      isLastChip,
      remove,
    };
  },
};
</script>
