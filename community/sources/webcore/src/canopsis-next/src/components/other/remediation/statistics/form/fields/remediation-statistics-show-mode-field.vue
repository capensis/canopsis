<template>
  <v-select
    v-field="value"
    :items="showingTypes"
    :label="$t('common.show')"
    hide-details
  />
</template>

<script>
import { REMEDIATION_STATISTICS_CHART_DATA_TYPE } from '@/constants';

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
          value: REMEDIATION_STATISTICS_CHART_DATA_TYPE.percent,
          text: this.$tc('common.percent'),
        },
        {
          value: REMEDIATION_STATISTICS_CHART_DATA_TYPE.remediation,
          text: this.$t('remediation.statistic.remediation'),
        },
      ];
    },
  },
};
</script>
