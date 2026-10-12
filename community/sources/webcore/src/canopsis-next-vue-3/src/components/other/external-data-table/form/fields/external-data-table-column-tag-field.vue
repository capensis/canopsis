<template>
  <c-select-chip
    :value="value"
    :items="items"
    :disabled="disabled"
    :color="activeChip?.color"
    class="px-2"
    text-color="white"
    bottom
    nudge-bottom
    @update:value="updateModel"
  />
</template>

<script>
import { computed } from 'vue';

import {
  EXTERNAL_DATA_TABLE_COLUMN_TAGS,
  EXTERNAL_DATA_TABLE_COLUMN_TYPES_COLORS,
} from '@/constants';

import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';

export default {
  emits: ['update:value'],
  props: {
    value: {
      type: Number,
      default: EXTERNAL_DATA_TABLE_COLUMN_TAGS.noTag,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateModel } = useModelField(props, emit);

    const items = computed(() => Object.values(EXTERNAL_DATA_TABLE_COLUMN_TAGS).map(value => ({
        value,
        text: t(`externalData.tableColumnTypes.${value}`),
        color: EXTERNAL_DATA_TABLE_COLUMN_TYPES_COLORS[value],
      })),
    );

    const activeChip = computed(() => items.value.find(({ value }) => value === props.value));

    return {
      activeChip,
      items,
      updateModel,
    };
  },
};
</script>
