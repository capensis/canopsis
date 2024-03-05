<template>
  <v-select
    v-field="value"
    :items="showingTypes"
    :label="$t('common.show')"
    hide-details
  />
</template>

<script>
import { KPI_SLI_GRAPH_DATA_TYPE } from '@/constants';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      required: true,
    },
  },
  computed: {
    showingTypes() {
      return [
        {
          value: KPI_SLI_GRAPH_DATA_TYPE.percent,
          text: this.$tc('common.percent'),
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
