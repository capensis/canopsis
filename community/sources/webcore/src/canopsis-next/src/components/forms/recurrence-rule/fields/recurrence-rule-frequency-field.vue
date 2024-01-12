<template>
  <v-select
    v-field="value"
    :label="$t('recurrenceRule.freq')"
    :items="frequencies"
  />
</template>

<script>
import { RRule } from 'rrule';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Number,
      required: false,
    },
  },
  computed: {
    frequencies() {
      return [
        { text: this.$t('recurrenceRule.none'), value: null },
        { text: this.$t('common.timeFrequencies.hourly'), value: RRule.HOURLY },
        { text: this.$t('common.timeFrequencies.daily'), value: RRule.DAILY },
        { text: this.$t('common.timeFrequencies.weekly'), value: RRule.WEEKLY },
        { text: this.$t('common.timeFrequencies.monthly'), value: RRule.MONTHLY },
        { text: this.$t('common.timeFrequencies.yearly'), value: RRule.YEARLY },
      ];
    },
  },
};
</script>
