<template>
  <v-row density="default">
    <v-col cols="6">
      <c-select
        :model-value="onSuccess"
        :label="$t('eventFilter.onSuccess')"
        :items="eventFilterAfterTypes"
        @update:model-value="$emit('update:on-success', $event)"
      />
    </v-col>
    <v-col cols="6">
      <c-select
        :model-value="onFailure"
        :label="$t('eventFilter.onFailure')"
        :items="eventFilterAfterTypes"
        @update:model-value="$emit('update:on-failure', $event)"
      />
    </v-col>
  </v-row>
</template>

<script>
import { computed } from 'vue';

import { EVENT_FILTER_ENRICHMENT_AFTER_TYPES } from '@/constants';

export default {
  emits: ['update:on-success', 'update:on-failure'],
  props: {
    onSuccess: {
      type: String,
      default: undefined,
    },
    onFailure: {
      type: String,
      default: undefined,
    },
  },
  setup() {
    const eventFilterAfterTypes = computed(() => Object.values(EVENT_FILTER_ENRICHMENT_AFTER_TYPES),
    );

    return {
      eventFilterAfterTypes,
    };
  },
};
</script>
