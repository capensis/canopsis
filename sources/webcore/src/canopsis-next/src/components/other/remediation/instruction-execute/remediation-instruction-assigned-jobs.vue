<template lang="pug">
  div
    v-layout(row)
      span.subheading {{ $t('remediationInstructionExecute.jobs.title') }}
    v-layout(column)
      v-data-table.jobs-assigned(
        :items="jobs",
        hide-actions
      )
        template(slot="headers", slot-scope="props")
          td
          td.text-xs-center {{ $t('remediationInstructionExecute.jobs.launchedAt') }}
          td.text-xs-center {{ $t('remediationInstructionExecute.jobs.completedAt') }}
        template(slot="items", slot-scope="props")
          td.pa-0
            v-btn.primary(
              round,
              small,
              block,
              @click="executeJob(props.item)"
            ) {{ props.item.name }}
          td.text-xs-center {{ props.item.launched_at || '-' }}
          td.text-xs-center {{ props.item.completed_at || '-' }}
</template>

<script>
export default {
  props: {
    jobs: {
      type: Array,
      default: () => [],
    },
    executionId: {
      type: String,
      required: true,
    },
    operationId: {
      type: [Number, String],
      required: true,
    },
  },
  methods: {
    executeJob() {},
  },
};
</script>

<style lang="scss">
  .jobs-assigned {
    tr {
      border-bottom: none !important;
    }

    tbody tr:hover {
      background: unset !important;
    }
  }
</style>
