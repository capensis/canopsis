<template lang="pug">
  c-advanced-data-table(
    :headers="headers",
    :items="failures",
    :loading="pending",
    :total-items="totalItems",
    :pagination="pagination",
    :is-expandable-item="hasEvent"
    expand,
    advanced-pagination,
    @update:pagination="$emit('update:pagination', $event)"
  )
    template(#status="{ item }")
      c-circle-badge.error.text-uppercase(v-if="item.unread") {{ $t('common.new') }}
    template(#t="{ item }") {{ item.t | date }}
    template(#type="{ item }") {{ $t(`eventFilter.failureTypes.${item.type}`) }}
    template(#message="{ item }")
      span.pre-wrap {{ item.message }}
    template(#expand="{ item }")
      event-filter-failures-list-expand-panel(:failure="item")
</template>

<script>
import EventFilterFailuresListExpandPanel from './event-filter-failures-list-expand-panel.vue';

export default {
  components: { EventFilterFailuresListExpandPanel },
  props: {
    failures: {
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
        { text: this.$t('common.timestamp'), value: 't' },
        { text: this.$t('common.type'), value: 'type' },
        { text: this.$t('common.message'), value: 'message' },
      ];
    },
  },
  methods: {
    hasEvent(failure) {
      return !!failure.event;
    },
  },
};
</script>
