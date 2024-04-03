<template>
  <c-advanced-data-table
    :items="steps"
    :headers="headers"
    hide-actions
  >
    <template #icon="{ item }">
      {{ item.completed_at | date('long', '-') }}
    </template>
    <template #completed_at="{ item }">
      {{ item.completed_at | date('long', '-') }}
    </template>
  </c-advanced-data-table>
</template>

<script>
export default {
  props: {
    steps: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    headers() {
      return [
        { text: '', value: 'icon', sortable: false },
        { text: this.$t('common.step'), value: 'name', sortable: false },
        { text: this.$t('common.output'), value: 'fail_reason', sortable: false },
        { text: this.$t('remediation.instructionExecute.jobs.completedAt'), value: 'completed_at', sortable: false },
      ];
    },
  },
};
</script>
