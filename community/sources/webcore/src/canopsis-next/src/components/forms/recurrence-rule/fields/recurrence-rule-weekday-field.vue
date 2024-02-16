<template>
  <v-layout
    v-if="chips"
    class="recurrence-rule-weekday-field"
    column
  >
    <v-subheader class="recurrence-rule-weekday-field__header pl-0">
      {{ $t('recurrenceRule.byweekday') }}
    </v-subheader>
    <v-chip-group
      v-field="value"
      active-class="elevation-2 grey lighten-1"
      multiple
    >
      <v-chip
        v-for="weekDay in weekDays"
        :key="weekDay.value"
        :value="weekDay.value"
      >
        {{ weekDay.text }}
      </v-chip>
    </v-chip-group>
  </v-layout>
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

<style lang="scss">
.recurrence-rule-weekday-field {
  &__header {
    height: unset;
  }
}
</style>
