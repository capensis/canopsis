<template lang="pug">
  v-data-iterator(:items="exceptions", hide-actions)
    template(#item="{ item, index }")
      v-layout(column)
        v-layout(row, align-center)
          h4 {{ item.name }}
          c-action-btn(type="delete", @click="removeItemFromArray(index)")
        c-advanced-data-table(:items="item.exdates", :headers="exdatesHeaders")
          template(#begin="{ item }") {{ item.begin | date }}
          template(#end="{ item }") {{ item.end | date }}
          template(#type="{ item }") {{ item.type.name }}
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
