<template>
  <v-data-iterator
    :items="exceptions"
    hide-default-footer
  >
    <template #item="{ item, index }">
      <v-layout column>
        <v-layout align-center>
          <h4>{{ item.name }}</h4>
          <c-action-btn
            type="delete"
            @click="removeItem(index)"
          />
        </v-layout>
        <c-advanced-data-table
          :items="item.exdates"
          :headers="exdatesHeaders"
        >
          <template #begin="{ item: exdate }">
            {{ exdate.begin | date }}
          </template>
          <template #end="{ item: exdate }">
            {{ exdate.end | date }}
          </template>
          <template #type="{ item: exdate }">
            {{ exdate.type.name }}
          </template>
        </c-advanced-data-table>
      </v-layout>
    </template>
  </v-data-iterator>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';
import { useArrayModelField } from '@/hooks/form/array-model-field';

export default {
  model: {
    prop: 'exceptions',
    event: 'input',
  },
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
      removeItem: removeItemFromArray,
    };
  },
};
</script>
