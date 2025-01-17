<template>
  <advanced-search-chip :items="" />
</template>
<script>
import AdvancedSearchChip from './advanced-search-chip.vue';

export default {
  props: {
    value: {
      type: String,
      default: '',
    },
  },
};
</script>
