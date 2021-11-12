<template lang="pug">
  v-autocomplete(
    v-field="value",
    :items="filters",
    :label="label || $t('common.filters')",
    :loading="filtersPending",
    :name="name",
    item-text="name",
    item-value="_id",
    clearable
  )
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
      default: 'role',
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      pending: false,
      items: [],
    };
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    fetchList() {
      if (!this.filtersPending || !this.filters.length) {
        this.fetchFiltersList({ params: { limit: MAX_LIMIT } });
      }
    },
  },
};
</script>
