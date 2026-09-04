<template>
  <v-data-table
    :items="items"
    :headers="headers"
    :items-per-page="itemsPerPage"
    hide-default-header
  />
</template>

<script>
import { computed } from 'vue';

import { PAGINATION_LIMIT } from '@/config';

export default {
  props: {
    names: {
      type: Array,
      default: () => [],
    },
  },

  setup(props) {
    const headers = [{ title: '', key: 'name', sortable: false }];
    const items = computed(() => props.names.map(name => ({ name })));

    return {
      headers,
      items,
      itemsPerPage: PAGINATION_LIMIT,
    };
  },
};
</script>
