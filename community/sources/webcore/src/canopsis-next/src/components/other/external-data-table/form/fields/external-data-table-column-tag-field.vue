<template>
  <c-select-chip
    v-field="value"
    :items="items"
    :value="value"
    :disabled="disabled"
    :color="activeChip?.color"
    class="px-2"
    text-color="white"
    bottom
    nudge-bottom
  />
</template>

<script>
import { computed } from 'vue';

import { EXTERNAL_DATA_TABLE_COLUMN_TAGS, EXTERNAL_DATA_TABLE_COLUMN_TYPES_COLORS } from '@/constants';

import { useI18n } from '@/hooks/i18n';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Number,
      default: EXTERNAL_DATA_TABLE_COLUMN_TAGS.noType,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const { t } = useI18n();

    const items = computed(() => Object.values(EXTERNAL_DATA_TABLE_COLUMN_TAGS).map(value => ({
      value,
      text: t(`externalData.tableColumnTypes.${value}`),
      color: EXTERNAL_DATA_TABLE_COLUMN_TYPES_COLORS[value],
    })));

    const activeChip = computed(() => items.value.find(({ value }) => value === props.value));

    return {
      activeChip,
      items,
    };
  },
};
</script>
