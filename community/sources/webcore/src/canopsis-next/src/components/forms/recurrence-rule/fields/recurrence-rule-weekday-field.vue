<template>
  <v-chip-group
    v-if="chips"
    v-field="value"
    :items="weekDays"
    :label="$t('recurrenceRule.byweekday')"
    multiple="multiple"
  />
  <v-select
    v-else
    v-field="value"
    :label="$t('recurrenceRule.wkst')"
    :items="weekDays"
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
      type: [Array, Number, String],
      default: () => [],
    },
    chips: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    weekDays() {
      return [
        { text: this.$t('common.weekDays.monday'), value: RRule.MO.weekday },
        { text: this.$t('common.weekDays.tuesday'), value: RRule.TU.weekday },
        { text: this.$t('common.weekDays.wednesday'), value: RRule.WE.weekday },
        { text: this.$t('common.weekDays.thursday'), value: RRule.TH.weekday },
        { text: this.$t('common.weekDays.friday'), value: RRule.FR.weekday },
        { text: this.$t('common.weekDays.saturday'), value: RRule.SA.weekday },
        { text: this.$t('common.weekDays.sunday'), value: RRule.SU.weekday },
      ];
    },
  },
};
</script>
