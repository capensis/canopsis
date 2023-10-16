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
            @click="removeItemFromArray(index)"
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
import { formArrayMixin } from '@/mixins/form';

export default {
  mixins: [formArrayMixin],
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
  computed: {
    exdatesHeaders() {
      return [
        { value: 'begin', text: this.$t('common.start'), sortable: false },
        { value: 'end', text: this.$t('common.end'), sortable: false },
        { value: 'type', text: this.$t('common.type'), sortable: false },
      ];
    },
  },
};
</script>
