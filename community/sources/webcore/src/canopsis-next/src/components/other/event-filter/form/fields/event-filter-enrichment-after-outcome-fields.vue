<template>
  <v-layout>
    <v-select
      v-field="form.config.on_success"
      :label="$t('eventFilter.onSuccess')"
      :items="eventFilterAfterTypes"
      class="mr-3"
    />
    <v-select
      v-field="form.config.on_failure"
      :label="$t('eventFilter.onFailure')"
      :items="eventFilterAfterTypes"
      class="ml-3"
    />
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { EVENT_FILTER_ENRICHMENT_AFTER_TYPES } from '@/constants';

export default {
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
  setup() {
    const eventFilterAfterTypes = computed(() => Object.values(EVENT_FILTER_ENRICHMENT_AFTER_TYPES));

    return {
      eventFilterAfterTypes,
    };
  },
};
</script>
