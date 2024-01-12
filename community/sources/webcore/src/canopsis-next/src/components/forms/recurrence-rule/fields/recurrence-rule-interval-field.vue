<template>
  <v-layout align-center>
    <v-flex xs6>
      <c-number-field
        v-field="value.interval"
        :label="$t('recurrenceRule.repeatEvery')"
        :min="1"
        name="interval"
      />
    </v-flex>
    <v-flex
      class="pl-2"
      v-if="value.interval"
      xs6
    >
      {{ intervalTimeString }}
    </v-flex>
  </v-layout>
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
      type: Object,
      required: true,
    },
  },
  computed: {
    intervalTimeString() {
      const timeMessageKey = {
        [RRule.HOURLY]: 'common.times.hour',
        [RRule.DAILY]: 'common.times.day',
        [RRule.WEEKLY]: 'common.times.week',
        [RRule.MONTHLY]: 'common.times.month',
        [RRule.YEARLY]: 'common.times.year',
      }[this.value.freq];

      return this.$te(timeMessageKey) ? this.$tc(timeMessageKey, this.value.interval) : '';
    },
  },
};
</script>
