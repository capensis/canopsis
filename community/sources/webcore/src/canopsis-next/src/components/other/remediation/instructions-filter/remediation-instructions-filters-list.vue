<template lang="pug">
  div
    remediation-instructions-filters-item(
      v-for="(filter, index) in filters",
      v-field="filters[index]",
      :key="filter._id",
      :filters="filters",
      :editable="editable",
      :closable="closable",
      @remove="removeItemFromArray(index)"
    )
</template>

<script>
import { formArrayMixin } from '@/mixins/form';

import RemediationInstructionsFiltersItem from './remediation-instructions-filters-item.vue';

export default {
  components: { RemediationInstructionsFiltersItem },
  mixins: [formArrayMixin],
  model: {
    prop: 'filters',
    event: 'input',
  },
  props: {
    filters: {
      type: Array,
      default: () => [],
    },
    editable: {
      type: Boolean,
      default: false,
    },
    closable: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
