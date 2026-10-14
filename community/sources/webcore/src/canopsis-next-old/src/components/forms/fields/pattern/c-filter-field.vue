<template>
  <v-autocomplete
    v-field="value"
    :items="filters"
    :label="label || $t('common.filters')"
    :loading="filtersPending"
    :disabled="disabled"
    :name="name"
    item-text="name"
    item-value="_id"
    hide-details
    clearable
  />
</template>

<script>
import { MAX_LIMIT } from '@/constants';

import { entitiesFilterMixin } from '@/mixins/entities/filter';

export default {
  mixins: [entitiesFilterMixin],
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
