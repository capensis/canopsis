<template lang="pug">
  v-chip-group(
    v-field="value",
    :items="months",
    :label="$t('recurrenceRule.bymonth')",
    multiple
  )
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
