<template>
  <c-select
    :model-value="value"
    :items="showingTypes"
    :label="$t('common.show')"
    item-title="text"
    item-value="value"
    hide-details
    @update:model-value="$emit('update:value', $event)"
  />
</template>

<script>
import { KPI_SLI_GRAPH_DATA_TYPE } from '@/constants';

export default {
  props: {
    value: {
      type: String,
      required: true,
    },
  },

  emits: ['update:value'],

  computed: {
    showingTypes() {
      return [
        {
          value: KPI_SLI_GRAPH_DATA_TYPE.percent,
          text: this.$t('common.percent'),
        },
        {
          value: KPI_SLI_GRAPH_DATA_TYPE.time,
          text: this.$t('common.time'),
        },
      ];
    },
  },
};
</script>
