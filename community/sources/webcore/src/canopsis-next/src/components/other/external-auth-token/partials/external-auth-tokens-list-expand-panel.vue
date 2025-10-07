<template>
  <div class="secondary pa-3">
    <v-card>
      <v-card-text>
        <v-data-table
          :headers="headers"
          :items="items"
        >
          <template #item.date="{ item }">
            {{ item.date | date }}
          </template>
        </v-data-table>
      </v-card-text>
    </v-card>
  </div>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    record: {
      type: Object,
      default: () => ({}),
    },
  },
  setup() {
    const { t } = useI18n();

    const headers = computed(() => [
      { text: t('common.requestDate'), value: 'date' },
      { text: t('common.output'), value: 'output' },
    ]);

    return {
      headers,
    };
  },
};
</script>
