<template lang="pug">
  v-layout.declare-ticket-test-query(column)
    v-layout(row, align-center, justify-space-between)
      v-flex(xs10)
        v-text-field(:label="$t('alarm.alarmDisplayName')")
      v-btn.white--text(color="orange") {{ $t('declareTicket.runTest') }}
</template>

<script>
export default {
  props: {
    declareTicketsRule: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
