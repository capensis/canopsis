<template lang="pug">
  c-advanced-data-table(
    :headers="headers",
    :items="errors",
    :loading="pending",
    :total-items="totalItems",
    :pagination="pagination",
    expand,
    advanced-pagination,
    @update:pagination="$emit('update:pagination', $event)"
  )
    template(#status="{ item }")
      c-circle-badge.error.text-uppercase(v-if="item.new") {{ $t('common.new') }}
    template(#timestamp="{ item }") {{ item.timestamp | date }}
</template>

<script>
export default {
  props: {
    errors: {
      type: Array,
      default: () => [],
    },
    totalItems: {
      type: Number,
      required: false,
    },
    pending: {
      type: Boolean,
      default: false,
    },
    pagination: {
      type: Object,
      required: true,
    },
  },
  computed: {
    headers() {
      return [
        { value: 'status', sortable: false, width: 50 },
        { text: this.$t('common.timestamp'), value: 'timestamp' },
        { text: this.$t('common.type'), value: 'type' },
        { text: this.$t('common.message'), value: 'message' },
      ];
    },
  },
};
</script>
