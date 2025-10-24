<template>
  <c-advanced-data-table
    :headers="headers"
    :items="failures"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :is-expandable-item="hasEvent"
    expand
    advanced-pagination
    @update:options="$emit('update:options', $event)"
  >
    <template #status="{ item }">
      <c-circle-badge
        v-if="item.unread"
        class="error text-uppercase"
      >
        {{ $t('common.new') }}
      </c-circle-badge>
    </template>
    <template #t="{ item }">
      {{ item.t | date }}
    </template>
    <template #type="{ item }">
      {{ $t(`eventFilter.failureTypes.${item.type}`) }}
    </template>
    <template #message="{ item }">
      <span class="pre-wrap">{{ item.message }}</span>
    </template>
    <template #expand="{ item }">
      <event-filter-failures-list-expand-panel :failure="item" />
    </template>
  </c-advanced-data-table>
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
    options: {
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
