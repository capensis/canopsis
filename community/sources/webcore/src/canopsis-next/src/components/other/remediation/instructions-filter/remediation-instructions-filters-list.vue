<template>
  <div>
    <remediation-instructions-filters-item
      v-field="filters[index]"
      v-for="(filter, index) in filters"
      :key="filter._id"
      :filters="filters"
      :editable="editable"
      :closable="closable"
      @remove="removeItemFromArray(index)"
    />
  </div>
</template>

<script>
import { formArrayMixin } from '@/mixins/form';

import RemediationInstructionsFiltersItem from './partials/remediation-instructions-filters-item.vue';

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
