<template>
  <v-autocomplete
    :model-value="value"
    :items="filters"
    :label="label || $t('common.filters')"
    :loading="filtersPending"
    :disabled="disabled"
    :name="name"
    item-title="name"
    item-value="_id"
    hide-details
    clearable
    @update:model-value="$emit('update:value', $event)"
  />
</template>

<script>
import { MAX_LIMIT } from '@/constants';

import { useFilter } from '@/hooks/store/modules/filter';

export default {
  props: {
    value: {
      type: [Object, String],
      required: false,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'filter',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  emits: ['update:value'],
  setup() {
    const { filters, filtersPending, fetchFiltersList } = useFilter();

    return {
      filters,
      filtersPending,
      fetchFiltersList,
    };
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    fetchList() {
      if (!this.filtersPending) {
        this.fetchFiltersList({ params: { limit: MAX_LIMIT } });
      }
    },
  },
};
</script>
