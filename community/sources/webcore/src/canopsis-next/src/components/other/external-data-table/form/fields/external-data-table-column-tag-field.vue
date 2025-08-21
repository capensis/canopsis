<template>
  <v-menu :disabled="disabled" bottom nudge-bottom>
    <template #activator="{ on }">
      <c-chip
        :color="activeChip?.color"
        class="px-2"
        text-color="white"
        v-on="on"
      >
        {{ activeChip?.text }}
      </c-chip>
    </template>
    <v-list>
      <v-list-item
        v-for="item in items"
        :key="item.value"
        :input-value="item.value === value"
        @click="select(item.value)"
      >
        <v-list-item-title>{{ item.text }}</v-list-item-title>
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script>
import { computed } from 'vue';

import { EXTERNAL_DATA_TABLE_COLUMN_TAGS, EXTERNAL_DATA_TABLE_COLUMN_TYPES_COLORS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModelField } from '@/hooks/form/model-field';

export default {
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
  setup(props, { emit }) {
    const { updateModel } = useModelField(props, emit);
    const { t } = useI18n();

    const items = computed(() => Object.values(EXTERNAL_DATA_TABLE_COLUMN_TAGS).map(value => ({
      value,
      text: t(`externalData.tableColumnTypes.${value}`),
      color: EXTERNAL_DATA_TABLE_COLUMN_TYPES_COLORS[value],
    })));

    const activeChip = computed(() => items.value.find(({ value }) => value === props.value));

    const select = updateModel;

    return {
      activeChip,
      items,
      select,
    };
  },
};
</script>
