<template>
  <v-layout
    class="recurrence-rule-month-field"
    column
  >
    <v-subheader class="recurrence-rule-month-field__header pl-0">
      {{ $t('recurrenceRule.bymonth') }}
    </v-subheader>
    <v-chip-group
      v-field="value"
      active-class="elevation-2 grey lighten-1"
      multiple
      column
    >
      <v-chip
        v-for="month in months"
        :key="month.value"
        :value="month.value"
      >
        {{ month.text }}
      </v-chip>
    </v-chip-group>
  </v-layout>
</template>

<script>
export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    months() {
      return [
        { text: this.$t('common.months.january'), value: 1 },
        { text: this.$t('common.months.february'), value: 2 },
        { text: this.$t('common.months.march'), value: 3 },
        { text: this.$t('common.months.april'), value: 4 },
        { text: this.$t('common.months.may'), value: 5 },
        { text: this.$t('common.months.june'), value: 6 },
        { text: this.$t('common.months.july'), value: 7 },
        { text: this.$t('common.months.august'), value: 8 },
        { text: this.$t('common.months.september'), value: 9 },
        { text: this.$t('common.months.october'), value: 10 },
        { text: this.$t('common.months.november'), value: 11 },
        { text: this.$t('common.months.december'), value: 12 },
      ];
    },
  },
};
</script>

<style lang="scss">
.recurrence-rule-month-field {
  &__header {
    height: unset;
  }
}
</style>
