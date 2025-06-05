<template>
  <v-card-text>
    <div class="text-center pa-4">
      <h3 class="mb-4">
        {{ $t('notifications.tabs.eventFilterFailures') }}
      </h3>
      <c-advanced-data-table
        :headers="headers"
        :items="items"
        :loading="pending"
        :total-items="meta.total_count"
        :options="options"
        advanced-pagination
        @update:options="$emit('update:options', $event)"
      />
    </div>
  </v-card-text>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    items: {
      type: Array,
      default: () => [],
    },
    pending: {
      type: Boolean,
      default: false,
    },
    meta: {
      type: Object,
      default: () => {},
    },
    options: {
      type: Object,
      default: () => {},
    },
  },
  setup() {
    const { t } = useI18n();

    // Computed
    const headers = computed(() => [
      { text: t('common.name'), value: 'name' },
      { text: t('common.description'), value: 'description' },
      { text: t('common.created'), value: 'created' },
    ]);

    return {
      headers,
    };
  },
};
</script>
