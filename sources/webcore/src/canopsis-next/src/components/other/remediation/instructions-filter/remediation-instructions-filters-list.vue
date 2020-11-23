<template lang="pug">
  div
    remediation-instructions-filters-item(
      v-for="(filter, index) in filters",
      v-field="filters[index]",
      :key="filter._id",
      :filters="filters",
      @remove="remove"
    )
</template>

<script>
import formBaseMixin from '@/mixins/form/base';

import RemediationInstructionsFiltersItem from './remediation-instructions-filters-item.vue';

export default {
  components: { RemediationInstructionsFiltersItem },
  mixins: [formBaseMixin],
  model: {
    prop: 'filters',
    event: 'input',
  },
  props: {
    filters: {
      type: Array,
      default: () => [],
    },
  },
  methods: {
    remove(filter) {
      this.updateModel(this.filters.filter(item => item._id !== filter._id));
    },
  },
};
</script>
