<template lang="pug">
  v-layout(column)
    v-card.my-2.py-2.px-3(v-for="(column, index) in columns", :key="index")
      column-field(
        v-field="columns[index]",
        :with-color-indicator="withColorIndicator",
        :with-html="withHtml",
        :with-template="withTemplate",
        :name="`${index}`",
        @up="up(index)",
        @down="down(index)",
        @remove="removeItemFromArray(index)"
      )
    v-btn.ml-0(color="primary", @click.prevent="add") {{ $t('common.add') }}
</template>

<script>
import { formArrayMixin, formValidationHeaderMixin } from '@/mixins/form';

import ColumnField from './column-field.vue';

export default {
  inject: ['$validator'],
  components: { ColumnField },
  mixins: [
    formArrayMixin,
    formValidationHeaderMixin,
  ],
  model: {
    prop: 'columns',
    event: 'input',
  },
  props: {
    columns: {
      type: [Array, Object],
      default: () => [],
    },
    withTemplate: {
      type: Boolean,
      default: false,
    },
    withHtml: {
      type: Boolean,
      default: false,
    },
    withColorIndicator: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    add() {
      const column = { label: '', value: '' };

      if (this.withHtml) {
        column.isHtml = false;
      }

      this.addItemIntoArray(column);
    },

    up(index) {
      if (index > 0) {
        const columns = [...this.columns];
        const temp = columns[index];

        columns[index] = columns[index - 1];
        columns[index - 1] = temp;

        this.updateModel(columns);
      }
    },

    down(index) {
      if (index < this.columns.length - 1) {
        const columns = [...this.columns];
        const temp = columns[index];

        columns[index] = columns[index + 1];
        columns[index + 1] = temp;

        this.updateModel(columns);
      }
    },
  },
};
</script>
