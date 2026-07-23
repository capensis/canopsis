<template>
  <v-data-iterator :items="exceptions" hide-default-footer>
    <template #item="{ item, index }">
      <c-layout column>
        <c-layout align="center">
          <h4>{{ item.name }}</h4>
          <c-action-btn type="delete" @click="removeItem(index)" />
        </c-layout>
        <c-advanced-data-table :items="item.exdates" :headers="exdatesHeaders">
          <template #begin="{ item: exdate }">
            {{ convertDateToString(exdate.begin) }}
          </template>
          <template #end="{ item: exdate }">
            {{ convertDateToString(exdate.end) }}
          </template>
          <template #type="{ item: exdate }">
            {{ exdate.type.name }}
          </template>
        </c-advanced-data-table>
      </c-layout>
    </template>
  </v-data-iterator>
</template>

<script>
import { computed } from 'vue';

import { convertDateToString } from '@/helpers/date/date';

import { useI18n } from '@/hooks/i18n';
import { useArrayModelField } from '@/hooks/form/array-model-field';

export default {
  emits: ['update:exceptions'],
  props: {
    exceptions: {
      type: Array,
      default: () => [],
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();

    const exdatesHeaders = computed(() => [
      { value: 'begin', text: t('common.start'), sortable: false },
      { value: 'end', text: t('common.end'), sortable: false },
      { value: 'type', text: t('common.type'), sortable: false },
    ]);

    const { removeItemFromArray } = useArrayModelField(props, emit);

    return {
      exdatesHeaders,
      convertDateToString,
      removeItem: removeItemFromArray,
    };
  },
};
</script>
